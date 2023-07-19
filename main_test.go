package main

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStep1Invalid(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "./testInputs/step1/invalid.json")
	err := cmd.Run()
	assert.Error(t, err, "exit status 1")
}

func TestStep1Valid(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "./testInputs/step1/valid.json")
	err := cmd.Run()
	assert.NoError(t, err)
}

func TestStep2Invalid(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "./testInputs/step2/invalid.json")
	err := cmd.Run()
	assert.Error(t, err, "exit status 1")
}

func TestStep2Invalid2(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "./testInputs/step2/invalid2.json")
	err := cmd.Run()
	assert.Error(t, err, "exit status 1")
}

func TestStep2Valid(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "./testInputs/step2/valid.json")
	err := cmd.Run()
	assert.NoError(t, err)
}

func TestStep2Valid2(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "./testInputs/step2/valid2.json")
	err := cmd.Run()
	assert.NoError(t, err)
}

func TestStep3Invalid(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "./testInputs/step3/invalid.json")
	err := cmd.Run()
	assert.Error(t, err, "exit status 1")
}

func TestStep3Valid(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "./testInputs/step3/valid.json")
	err := cmd.Run()
	assert.NoError(t, err)
}

func TestStep4Invalid(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "./testInputs/step4/invalid.json")
	err := cmd.Run()
	assert.Error(t, err, "exit status 1")
}

func TestStep4Valid(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "./testInputs/step4/valid.json")
	err := cmd.Run()
	assert.NoError(t, err)
}

func TestStep4Valid2(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "./testInputs/step4/valid2.json")
	err := cmd.Run()
	assert.NoError(t, err)
}
