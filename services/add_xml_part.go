package services

import "mime/multipart"

func AddXMLPart(mw *multipart.Writer, content string) error {
	part, err := mw.CreatePart(map[string][]string{
		"Content-Type": {"application/xml; charset=utf-8"},
	})
	if err != nil {
		return err
	}
	_, err = part.Write([]byte(content))

	if err != nil {
		return err
	}

	return nil
}
