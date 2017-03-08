package openface

import (
    "os/exec"
)

func Train(imgDir string, alignDir string, featureDir string) error {
    err := AlignImages(imgDir, alignDir)
    if err != nil {
       return err
    }

    err = GenReps(alignDir, featureDir)
    if err != nil {
        return err
    }

    err = CreatePickle(featureDir)
    if err != nil {
        return err
    }
}

func AlignImages(imgDir string, alignDir string) error {
    script := "/root/openface/scripts/align.sh"
    cmd := exec.Command(script, imgDir, alignDir)
    err = cmd.Run()
    if err != nil {
        return err
    }
    return nil
}

func GenReps(alignDir string, featureDir string) error {
    script := "/root/openface/batch-represent/main.lua"
    cmd := exec.Command(script, "-outDir", featureDir, "-data", alignDir)
    err = cmd.Run()
    if err != nil {
        return err
    }
    return nil
}

func CreatePickle(featureDir string) error {
    script := "/root/openface/scripts/classifier.py"
    cmd := exec.Command(script, "train", featureDir)
    err := cmd.Run()
    if err != nil {
        return err
    }
    return nil
}
