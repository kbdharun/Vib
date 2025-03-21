package core

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/mitchellh/mapstructure"
	"github.com/vanilla-os/vib/api"
)

// Compile and build the recipe using the specified runtime
func CompileRecipe(recipePath string, runtime string, isRoot bool, origGid int, origUid int) error {
	recipe, err := BuildRecipe(recipePath)
	if err != nil {
		return err
	}

	syscall.Setegid(0)
	syscall.Seteuid(0)
	switch runtime {
	case "docker":
		err = compileDocker(recipe, origGid, origUid)
		if err != nil {
			return err
		}
	case "podman":
		err = compilePodman(recipe, origGid, origUid)
		if err != nil {
			return err
		}
	case "buildah":
		return fmt.Errorf("buildah not implemented yet")
	default:
		return fmt.Errorf("no runtime specified and the prometheus library is not implemented yet")
	}
	syscall.Setegid(origGid)
	syscall.Seteuid(origUid)

	for _, finalizeInterface := range recipe.Finalize {
		var module Finalize

		err := mapstructure.Decode(finalizeInterface, &module)
		if err != nil {
			return err
		}
		err = LoadFinalizePlugin(module.Type, finalizeInterface, &recipe, runtime, isRoot, origGid, origUid)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Image %s built successfully using %s\n", recipe.Id, runtime)

	return nil
}

// Build an OCI image using the specified recipe through Docker
func compileDocker(recipe api.Recipe, gid int, uid int) error {
	docker, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	cmd := exec.Command(
		docker, "build",
		"-t", fmt.Sprintf("localhost/%s", recipe.Id),
		"-f", recipe.Containerfile,
		".",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = recipe.ParentPath

	return cmd.Run()
}

// Build an OCI image using the specified recipe through Podman
func compilePodman(recipe api.Recipe, gid int, uid int) error {
	podman, err := exec.LookPath("podman")
	if err != nil {
		return err
	}

	cmd := exec.Command(
		podman, "build",
		"-t", fmt.Sprintf("localhost/%s", recipe.Id),
		"-f", recipe.Containerfile,
		".",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = recipe.ParentPath

	return cmd.Run()
}
