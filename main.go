package main

import (
    "fmt"
    "gioui.org/app"
    "gioui.org/f32"
    "gioui.org/io/system"
    "gioui.org/layout"
    "gioui.org/op"
    "gioui.org/op/paint"
    "gioui.org/unit"

    "image"
    //"image/color"
    _ "image/jpeg"
    "log"
    "os"
)

func main() {
    var rawImg image.Image
    //getRandomPath(getRandomPage())
    //getRandomPath()

    rawImg = getWall()
    scalingFactor := getScalingFactor(rawImg.Bounds().Max.X, rawImg.Bounds().Max.Y)
    sizeX := float32(rawImg.Bounds().Max.X) * scalingFactor
    sizeY := float32(rawImg.Bounds().Max.Y) * scalingFactor

    go func() {
        w := app.NewWindow(app.Title("Drawing Ref Wallhaven"), app.Size(unit.Dp(sizeX), unit.Dp(sizeY)))
        if err := loop(w, &rawImg, scalingFactor); err != nil {
            log.Fatal(err)
        }
        os.Exit(0)
    }()
    app.Main()
}

func getImage() image.Image {
    f, err := os.Open("test.jpg")
    if err != nil {
        fmt.Println(err)
    }
    defer f.Close()

    img, _, err := image.Decode(f)
    if err != nil {
        fmt.Println(err)
    }
    return img
}

func loop(w *app.Window, rawImg *image.Image, scalingFactor float32) error {
    var ops op.Ops
    for {
        e := <-w.Events()
        switch e := e.(type) {
        case system.DestroyEvent:
            return e.Err
        case system.FrameEvent:

            gtx := layout.NewContext(&ops, e)

            layout.UniformInset(unit.Dp(0)).Layout(gtx,
                func(gtx layout.Context) layout.Dimensions {
                    //size := gtx.Constraints.Max

                    //rawImg := getWall()
                    img := paint.NewImageOp(*rawImg)
                    size := image.Point{X: img.Size().X, Y: img.Size().Y}
                    //op.Affine(
                    //    f32.Affine2D{}.Scale(f32.Point{}, f32.Pt(float32(size.X) / float32(img.Size().X),
                    //        float32(size.Y) / float32(img.Size().Y))),
                    //
                    //).Add(gtx.Ops)
                    scale := f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(scalingFactor, scalingFactor))
                    op.Affine(scale).Add(gtx.Ops)

                    img.Add(gtx.Ops)
                    paint.PaintOp{}.Add(gtx.Ops)
                    return layout.Dimensions{Size: size}
                },
            )
            e.Frame(gtx.Ops)
        }
    }
}

func getScalingFactor(x, y int) float32 {
    resolutionX := 1920 - 50
    resolutionY := 1080 - 100
    ratioX := float32(x) / float32(resolutionX)
    ratioY := float32(y) / float32(resolutionY)
    var scale float32
    if x < resolutionX && y < resolutionY {
        scale = 1
        return scale
    }
    if ratioX > ratioY {
        scale = 1 / ratioX
        return scale
    }
    if y < resolutionY {
        scale = 1
        return scale
    }
    scale = 1 / ratioY
    return scale
}
