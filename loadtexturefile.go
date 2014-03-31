// loadtexturefile.go

// Copyright 2012 Lightpoke. All rights reserved.
// This source code is subject to the terms and
// conditions defined in the "License.txt" file.

// Package texture is a fork from Azul3d.org exposing LoadTextureFileHandle.
// Instead of using LoadTextureFile or directly embedding with LoadTexture,
// LoadTextureFileHandle allows loading from an embedding packages file handle
// and resizing of textures that are too large.
//
// This should not be used in production. You should verify the size of the texture
// and resize them if needed prior to embedding.
package texture

import (
	"azul3d.org/v0/scene"
	"azul3d.org/v0/scene/renderer"
	"azul3d.org/v0/scene/texture"
	"azul3d.org/v0/thirdparty/resize"
	"image"
	"image/draw"
	"io"
	"log"
	"os"
)

// LoadTextureFile is the same as original azul3d.
func LoadTextureFile(n *scene.Node, filePath string) (texture.Type, error) {
	// Open texture file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return LoadTextureFileHandle(n, file), nil
}

// LoadTextureFileHandle loads and resizes texture using file handle.
func LoadTextureFileHandle(n *scene.Node, file io.ReadCloser) texture.Type {
	// Create texture object using no texture
	tex2 := texture.New()

	// Spawn goroutine to do texture loading
	go func() {
		defer file.Close()

		// Decode image
		srcImage, _, err := image.Decode(file)
		if err != nil {
			log.Printf("Load():", err)
		}

		// Check to see if we need to resize image because it is too large.
		maxSize := int(renderer.MaxTextureSize(n))
		sz := srcImage.Bounds().Size()
		if sz.X > maxSize || sz.Y > maxSize {
			// Cannot upload texture larger than max texture size, we must resize
			// it first.
			newWidth := sz.X
			if newWidth > maxSize {
				newWidth = maxSize
			}
			newHeight := sz.Y
			if newHeight > maxSize {
				newHeight = maxSize
			}

			srcImage = resize.Resample(srcImage, srcImage.Bounds(), newWidth, newHeight)
		}

		rgbaImage, ok := srcImage.(*image.RGBA)
		if !ok {
			// Convert image to RGBA
			b := srcImage.Bounds()
			rgbaImage = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
			draw.Draw(rgbaImage, rgbaImage.Bounds(), srcImage, b.Min, draw.Src)
		}

		// Assign image to texture
		tex2.SetImage(rgbaImage)

		// Tell renderer to load the texture
		renderer.LoadTexture(n, tex2)
	}()

	return tex2
}
