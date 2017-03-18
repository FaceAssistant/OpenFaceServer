package openface

import (
    "os/exec"
    "fa/s3util"
    "fmt"
    "os"
)

func Train(imgDir string, alignDir string, featureDir string, userId int) error {
    err := AlignImages(imgDir, alignDir)
    if err != nil {
       return err
    }

    err = GenReps(alignDir, featureDir)
    if err != nil {
        return err
    }

    err = ConcatFeatures(featureDir, userId)
    if err != nil {
        return err
    }

    err = CreatePickle(featureDir)
    if err != nil {
        return err
    }
    return nil
}

func AlignImages(imgDir string, alignDir string) error {
    script := "/root/openface/scripts/align.sh"
    cmd := exec.Command(script, imgDir, alignDir)
    err := cmd.Run()
    if err != nil {
        return err
    }
    return nil
}

func GenReps(alignDir string, featureDir string) error {
    script := "/root/openface/batch-represent/main.lua"
    cmd := exec.Command(script, "-outDir", featureDir, "-data", alignDir)
    err := cmd.Run()
    if err != nil {
        return err
    }

    return nil
}

func ConcatFeatures(featureDir string, userId int) error {
    labels, err := os.OpenFile(fmt.Sprintf("%s/labels.csv", featureDir), os.O_APPEND|os.O_WRONLY, 0600)
    defer labels.Close()
    if err != nil {
        return err
    }

    reps, err := os.OpenFile(fmt.Sprintf("%s/reps.csv", featureDir), os.O_APPEND|os.O_WRONLY, 0600)
    defer reps.Close()
    if err != nil {
        return err
    }

    err = s3util.GetFeature("labels.csv", userId, labels)
    if err != nil {
        return err
    }

    err = s3util.GetFeature("reps.csv", userId, reps)
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
