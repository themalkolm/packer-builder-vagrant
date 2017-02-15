package oneandone

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	set_image  sync.Once
	image_name string
	image_desc string
	test_image *Image
	image_serv *Server
)

const (
	img_freq = "WEEKLY"
	img_numb = 1
)

// Helper functions

func create_image(ser_id string) *Image {
	rand.Seed(time.Now().UnixNano())
	ri := rand.Intn(1000)
	image_name = fmt.Sprintf("TestServerImage_%d", ri)
	image_desc = fmt.Sprintf("TestServerImage_%d description", ri)
	req := ImageConfig{
		Name:        image_name,
		Description: image_desc,
		ServerId:    ser_id,
		Frequency:   img_freq,
		NumImages:   img_numb,
	}
	fmt.Printf("Creating image '%s'...\n", image_name)
	img_id, img, err := api.CreateImage(&req)
	if err != nil {
		fmt.Printf("Unable to create new image. Error: %s", err.Error())
		return nil
	}
	if img_id == "" || img.Id == "" {
		fmt.Printf("Unable to create image '%s'.", image_name)
		return nil
	}
	api.WaitForState(img, "ENABLED", 10, 90)
	return img
}

func setup_image() {
	_, image_serv, _ = create_test_server(true)
	api.WaitForState(image_serv, "POWERED_ON", 10, 180)
	test_image = create_image(image_serv.Id)
}

// /images tests

func TestCreateImage(t *testing.T) {
	set_image.Do(setup_image)

	if test_image == nil {
		t.Errorf("CreateImage failed.")
		return
	}
	if test_image.Id == "" {
		t.Errorf("Missing image ID.")
		time.Sleep(60 * time.Second)
	}
	if !strings.Contains(test_image.Name, image_name) {
		t.Errorf("Wrong image name.")
		time.Sleep(60 * time.Second)
	}
	if test_image.Description != image_desc {
		t.Errorf("Wrong image description.")
		time.Sleep(60 * time.Second)
	}
	if test_image.ServerId != image_serv.Id {
		t.Errorf("Wrong server ID in image '%s'.", test_image.Name)
		time.Sleep(60 * time.Second)
	}
	if test_image.Frequency != img_freq {
		t.Errorf("Wrong image frequency.")
		time.Sleep(60 * time.Second)
	}
	if test_image.NumImages != img_numb {
		t.Errorf("Wrong number of images in image '%s'.", test_image.Name)
		time.Sleep(60 * time.Second)
	}
}

func TestGetImage(t *testing.T) {
	set_image.Do(setup_image)

	fmt.Printf("Getting image '%s'...\n", test_image.Name)
	img, err := api.GetImage(test_image.Id)

	if err != nil {
		t.Errorf("GetImage failed. Error: " + err.Error())
		return
	}
	if img.Name != test_image.Name {
		t.Errorf("Wrong image name.")
	}
	if img.Type != test_image.Type {
		t.Errorf("Wrong image type.")
	}
	if *img.Architecture != *test_image.Architecture {
		t.Errorf("Wrong image architecture.")
	}
	if img.Description != test_image.Description {
		t.Errorf("Wrong image description.")
	}
	if img.ServerId != test_image.ServerId {
		t.Errorf("Wrong server ID in image '%s'.", test_image.Name)
	}
	if img.Frequency != test_image.Frequency {
		t.Errorf("Wrong image frequency.")
	}
	if img.NumImages != test_image.NumImages {
		t.Errorf("Wrong number of images in image '%s'.", test_image.Name)
	}
}

func TestListImages(t *testing.T) {
	set_image.Do(setup_image)
	fmt.Println("Listing all images...")

	imgs, err := api.ListImages()
	if err != nil {
		t.Errorf("ListImages failed. Error: " + err.Error())
	}
	if len(imgs) == 0 {
		t.Errorf("No image found.")
	}

	imgs, err = api.ListImages(1, 3, "name", "", "id,name,type")
	if err != nil {
		t.Errorf("ListImages with parameter options failed. Error: " + err.Error())
		return
	}
	if len(imgs) == 0 {
		t.Errorf("No image found.")
	}
	if len(imgs) > 3 {
		t.Errorf("Wrong number of objects per page.")
	}
	if imgs[0].Id == "" {
		t.Errorf("Filtering parameters failed.")
	}
	if imgs[0].Name == "" {
		t.Errorf("Filtering parameters failed.")
	}
	if imgs[0].State != "" {
		t.Errorf("Filtering parameters failed.")
	}
	if len(imgs) >= 2 && imgs[0].Name > imgs[1].Name {
		t.Errorf("Sorting parameters failed.")
	}

	imgs, err = api.ListImages(0, 0, "", test_image.Name, "")
	if err != nil {
		t.Errorf("ListImages with parameter options failed. Error: " + err.Error())
		return
	}
	if len(imgs) != 1 {
		t.Errorf("Search parameter failed.")
	}
	if imgs[0].Name != test_image.Name {
		t.Errorf("Search parameter failed.")
	}
}

func TestUpdateImage(t *testing.T) {
	set_image.Do(setup_image)

	fmt.Printf("Updating image '%s'...\n", test_image.Name)
	new_name := test_image.Name + " updated"
	new_desc := test_image.Name + " updated"

	img, err := api.UpdateImage(test_image.Id, new_name, new_desc, "ONCE")

	if err != nil {
		t.Errorf("UpdateImage failed. Error: " + err.Error())
		return
	}
	api.WaitForState(img, "ACTIVE", 10, 30)

	if img.Id != test_image.Id {
		t.Errorf("Wrong image ID.")
	}
	if img.Name != new_name {
		t.Errorf("Unable to update image '%s' name.", img.Name)
	}
	if img.Description != new_desc {
		t.Errorf("Unable to update image '%s' description.", img.Name)
	}
	if img.Frequency != "ONCE" {
		t.Errorf("Unable to update image '%s' frequency.", img.Name)
	}
}

func TestDeleteImage(t *testing.T) {
	set_image.Do(setup_image)

	fmt.Printf("Deleting image '%s'...\n", test_image.Name)
	img, err := api.DeleteImage(test_image.Id)

	if err != nil {
		t.Errorf("DeleteImage failed. Error: " + err.Error())
		return
	} else {
		api.WaitUntilDeleted(img)
	}

	img, _ = api.GetImage(test_image.Id)

	if img != nil {
		t.Errorf("Unable to delete the image.")
	} else {
		test_image = nil
	}
	// Delete test server
	api.DeleteServer(image_serv.Id, false)
	api.WaitUntilDeleted(image_serv)
	image_serv, _ = api.GetServer(image_serv.Id)
}
