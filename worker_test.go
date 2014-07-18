package ttygif

import (
	"testing"
)

func TestWorker(t *testing.T) {
	worker := NewWorker()
	worker.AddTargetFile("./test_data/test.png", "png")
	worker.AddTargetFile("./test_data/test.xwd", "xwd")

	images, err := worker.GetAllImages()
	if err != nil {
		t.Error(err)
	}
	if len(images) != 2 {
		t.Errorf("num of images: %d", len(images))
	}
}
