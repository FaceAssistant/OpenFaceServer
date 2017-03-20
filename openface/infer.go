package openface

import (
    "os/exec"
    "strings"
)

func Infer(classifier string, image string) ([]string, error) {
    script  :=  "/root/openface/scripts/classifier.py"
    arg0 := "infer"
    arg1 := classifier
    arg2 := image

    cmd := exec.Command(script, arg0, arg1, arg2)
    stdout, err := cmd.CombinedOutput()
    if err != nil {
        return nil, err
    }

    strs := strings.Split(strings.Trim(string(stdout), "\n"), ",")
    return strs, nil
}
