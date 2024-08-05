package utils

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/BeamStackProj/beamstack-cli/src/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/remotecommand"
)

func MigrateFilesToContainer(clientset *kubernetes.Clientset, params types.MigrationParams) error {
	config := GetKubeConfig()
	pod, err := clientset.CoreV1().Pods(params.Pod.Namespace).Get(context.TODO(), params.Pod.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	containers := pod.Spec.Containers
	if params.ContainerName == nil {
		params.ContainerName = &containers[0].Name
	}

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	srcFileInfo, err := os.Stat(params.SrcPath)
	if err != nil {
		return fmt.Errorf("error loading file info %s", err)
	}

	if srcFileInfo.IsDir() {
		err = filepath.Walk(params.SrcPath, func(file string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := tar.FileInfoHeader(fi, fi.Name())
			if err != nil {
				return err
			}
			relativePath := strings.TrimPrefix(file, params.SrcPath)
			relativePath = strings.TrimPrefix(relativePath, string(filepath.Separator))
			header.Name = filepath.ToSlash(relativePath)

			if err := tw.WriteHeader(header); err != nil {
				return err
			}
			if fi.Mode().IsRegular() {
				srcFile, err := os.Open(file)
				if err != nil {
					return fmt.Errorf("error opening file %s: %v", file, err)
				}
				defer srcFile.Close()

				if _, err := io.Copy(tw, srcFile); err != nil {
					return fmt.Errorf("error copying file %s: %v", file, err)
				}
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error walking source path %s: %v", params.SrcPath, err)
		}
	} else {
		// Handle single file
		srcFile, err := os.Open(params.SrcPath)
		if err != nil {
			return fmt.Errorf("error opening file %s", err)
		}
		defer srcFile.Close()

		header := &tar.Header{
			Name: filepath.Base(params.SrcPath),
			Mode: int64(srcFileInfo.Mode()),
			Size: srcFileInfo.Size(),
		}
		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("error writing header %s", err)
		}
		if _, err := io.Copy(tw, srcFile); err != nil {
			return fmt.Errorf("error copying file %s", err)
		}
	}

	if err := tw.Close(); err != nil {
		return fmt.Errorf("error closing tar writer: %v", err)
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(params.Pod.Name).
		Namespace(params.Pod.Namespace).
		SubResource("exec").
		Param("container", *params.ContainerName).
		Param("command", "/bin/sh").
		Param("command", "-c").
		Param("command", "mkdir -p "+params.DestPath+" && tar xvf - -C "+params.DestPath).
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "false")

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}

	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdin:  buf,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if err != nil {
		return err
	}

	return nil
}

func MigrateFilesFromContainer(clientset *kubernetes.Clientset, params types.MigrationParams) error {
	config := GetKubeConfig()

	pod, err := clientset.CoreV1().Pods(params.Pod.Namespace).Get(context.TODO(), params.Pod.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	containers := pod.Spec.Containers
	if params.ContainerName == nil {
		params.ContainerName = &containers[0].Name
	}

	buf := new(bytes.Buffer)

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(params.Pod.Name).
		Namespace(params.Pod.Namespace).
		SubResource("exec").
		Param("container", *params.ContainerName).
		Param("command", "/bin/sh").
		Param("command", "-c").
		Param("command", "tar cf - "+params.SrcPath).
		Param("stdin", "false").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "false")

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}

	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdout: buf,
		Stderr: os.Stderr,
	})
	if err != nil {
		return err
	}

	tr := tar.NewReader(buf)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		normalizedSrcPath := filepath.ToSlash(strings.TrimPrefix(params.SrcPath, "/"))
		normalizedHeaderName := filepath.ToSlash(strings.TrimPrefix(header.Name, "/"))

		relativePath := strings.TrimPrefix(normalizedHeaderName, normalizedSrcPath)
		target := filepath.Join(params.DestPath, relativePath)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}

	return nil
}
