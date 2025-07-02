package modelT

import (
	DTO "agenda-kaki-go/core/config/api/dto"
	"agenda-kaki-go/core/config/db/model"
	"agenda-kaki-go/core/config/namespace"
	"agenda-kaki-go/core/lib"
	FileBytes "agenda-kaki-go/core/lib/file_bytes"
	handler "agenda-kaki-go/core/test/handlers"
	"bytes"
	"fmt"
	"reflect"
)

type Client struct {
	Created      *model.Client
	Appointments []*Appointment
	X_Auth_Token string
}

func (u *Client) Set() error {
	if err := u.Create(200); err != nil {
		return err
	}
	if err := u.VerifyEmail(200); err != nil {
		return err
	}
	if err := u.Login(200); err != nil {
		return err
	}
	if err := u.UploadImages(200, map[string][]byte{
		"profile": FileBytes.PNG_FILE_1,
	}, nil); err != nil {
		return err
	}
	return nil
}

func (u *Client) Create(s int) error {
	pswd := "1SecurePswd!"
	if err := handler.NewHttpClient().
		Method("POST").
		URL("/client").
		ExpectedStatus(s).
		Send(DTO.CreateClient{
			Email:    lib.GenerateRandomEmail("client"),
			Name:     lib.GenerateRandomName("Client Name"),
			Surname:  lib.GenerateRandomName("Client Surname"),
			Password: pswd,
			Phone:    lib.GenerateRandomPhoneNumber(),
		}).ParseResponse(&u.Created).Error; err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	u.Created.Password = pswd
	return nil
}

func (u *Client) Update(s int, changes map[string]any) error {
	if err := handler.NewHttpClient().
		Method("PATCH").
		URL("/client/"+fmt.Sprintf("%v", u.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, u.X_Auth_Token).
		Send(changes).
		ParseResponse(&u.Created).
		Error; err != nil {
		return fmt.Errorf("failed to update client: %w", err)
	}
	if s > 200 && s < 300 {
		if err := ValidateUpdateChanges("Client", u.Created, changes); err != nil {
			return err
		}
	}

	return nil
}

func (u *Client) GetByEmail(s int) error {
	if err := handler.NewHttpClient().
		Method("GET").
		URL("/client/email/"+u.Created.Email).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, u.X_Auth_Token).
		Send(nil).
		ParseResponse(&u.Created).Error; err != nil {
		return fmt.Errorf("failed to get client by email: %w", err)
	}
	return nil
}

func (u *Client) Delete(s int) error {
	if err := handler.NewHttpClient().
		Method("DELETE").
		URL(fmt.Sprintf("/client/%v", u.Created.ID.String())).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, u.X_Auth_Token).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}
	return nil
}

func (u *Client) VerifyEmail(s int) error {
	if err := handler.NewHttpClient().
		Method("POST").
		URL(fmt.Sprintf("/client/verify-email/%v/%s", u.Created.Email, "12345")).
		ExpectedStatus(s).
		Header(namespace.HeadersKey.Auth, u.X_Auth_Token).
		Send(nil).
		Error; err != nil {
		return fmt.Errorf("failed to verify client email: %w", err)
	}
	return nil
}

func (u *Client) Login(s int) error {
	login := DTO.LoginClient{
		Email:    u.Created.Email,
		Password: "1SecurePswd!",
	}
	http := handler.NewHttpClient()
	if err := http.
		Method("POST").
		URL("/client/login").
		ExpectedStatus(s).
		Send(login).Error; err != nil {
		return fmt.Errorf("failed to login client: %w", err)
	}
	auth := http.ResHeaders[namespace.HeadersKey.Auth]
	if len(auth) == 0 {
		return fmt.Errorf("authorization header '%s' not found", namespace.HeadersKey.Auth)
	}
	u.X_Auth_Token = auth[0]
	if err := u.GetByEmail(200); err != nil {
		return fmt.Errorf("failed to get client by email after login: %w", err)
	}
	return nil
}

func ValidateUpdateChanges(modelName string, v any, changes map[string]any) error {
	mappy, err := lib.StructToMap(v)
	if err != nil {
		return fmt.Errorf("failed to convert %s struct to map: %w", modelName, err)
	}

	for key, expected := range changes {
		// Se o expected for struct, transforma em map
		if reflect.TypeOf(expected).Kind() == reflect.Struct {
			expected, err = lib.StructToMap(expected)
			if err != nil {
				return fmt.Errorf("failed to convert expected value for key '%s' to map: %w", key, err)
			}
		}

		actual := mappy[key]

		if !reflect.DeepEqual(actual, expected) {
			return fmt.Errorf("%s %s was not updated: expected '%#v' but got '%#v'", modelName, key, expected, actual)
		}
	}

	return nil
}

func (c *Client) UploadImages(status int, files map[string][]byte, x_auth_token *string) error {
	var fileMap = make(handler.Files)
	for field, content := range files {
		fileMap[field] = handler.MyFile{
			Name:    field + "_" + lib.GenerateRandomString(6) + ".png",
			Content: content,
		}
	}

	t, err := get_token(x_auth_token, &c.X_Auth_Token)
	if err != nil {
		return err
	}

	if err := handler.NewHttpClient().
		Method("PATCH").
		URL(fmt.Sprintf("/client/%s/design/images", c.Created.ID.String())).
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Auth, t).
		Send(fileMap).
		ParseResponse(&c.Created.Design.Images).
		Error; err != nil {
		return fmt.Errorf("failed to upload client images: %w", err)
	}

	return nil
}

func (c *Client) GetImage(status int, imageURL string, compareImgBytes *[]byte) error {
	if imageURL == "" {
		return fmt.Errorf("image URL cannot be empty")
	}
	http := handler.NewHttpClient()
	http.Method("GET")
	http.URL(imageURL)
	http.ExpectedStatus(status)
	http.Send(nil)
	// Compare the response bytes with the expected image bytes
	if compareImgBytes != nil {
		var response []byte
		http.ParseResponse(&response)
		if len(response) == 0 {
			return fmt.Errorf("received empty response for image (%s)", imageURL)
		} else if len(response) != len(*compareImgBytes) {
			return fmt.Errorf("image size mismatch for %s: expected %d bytes, got %d bytes", imageURL, len(*compareImgBytes), len(response))
		} else if !bytes.Equal(response, *compareImgBytes) {
			return fmt.Errorf("image content mismatch for %s", imageURL)
		}
	}
	return nil
}

func (c *Client) DeleteImages(status int, image_types []string, x_auth_token *string) error {
	if len(image_types) == 0 {
		return fmt.Errorf("no image types provided to delete")
	}

	t, err := get_token(x_auth_token, &c.X_Auth_Token)
	if err != nil {
		return err
	}

	http := handler.NewHttpClient()

	if err := http.
		Method("DELETE").
		ExpectedStatus(status).
		Header(namespace.HeadersKey.Auth, t).
		Error; err != nil {
		return fmt.Errorf("failed to prepare delete images request: %w", err)
	}

	base_url := fmt.Sprintf("/client/%s/design/images", c.Created.ID.String())
	for _, image_type := range image_types {
		image_url := base_url + "/" + image_type
		http.URL(image_url)
		http.Send(nil)
		http.ParseResponse(&c.Created.Design.Images)
		if http.Error != nil {
			return fmt.Errorf("failed to delete image %s: %w", image_type, http.Error)
		}
		url := c.Created.Design.Images.GetImageURL(image_type)
		if url != "" {
			return fmt.Errorf("image %s was not deleted successfully, expected empty URL but got %s", image_type, url)
		}
	}
	return nil
}
