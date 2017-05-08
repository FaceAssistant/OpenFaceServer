package openface

import (
    "os/exec"
    "fa/s3util"
    "fmt"
    "io"
    "os"
)

func Train(imgDir string, alignDir string, featureDir string, userId string, idToken string) error {
    err := AlignImages(imgDir, alignDir)
    if err != nil {
       return err
    }

    err = GenReps(alignDir, featureDir)
    if err != nil {
        return err
    }

    err = ConcatFeatures(featureDir, userId, idToken)
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

func ConcatFeatures(featureDir string, userId string, idToken string) error {
    labels, err := os.OpenFile(fmt.Sprintf("%s/labels.csv", featureDir), os.O_APPEND|os.O_RDWR, 0666)
    if err != nil {
        return err
    }
    defer labels.Close()

    reps, err := os.OpenFile(fmt.Sprintf("%s/reps.csv", featureDir), os.O_APPEND|os.O_RDWR, 0666)
    if err != nil {
        return err
    }
    defer reps.Close()

    l, err := os.Create(featureDir + "/labels")
    if err != nil {
        return err
    }
    defer l.Close()

    r, err := os.Create(featureDir + "/reps")
    if err != nil {
        return err
    }
    defer r.Close()

    _, err = io.Copy(l, labels)
    if err != nil {
        return err
    }

    _, err = io.Copy(r, reps)
    if err != nil {
        return err
    }

    err = s3util.GetFeature("labels.csv", userId, idToken, labels)
    if err != nil {
        return err
    }

    err = s3util.GetFeature("reps.csv", userId, idToken, reps)
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
