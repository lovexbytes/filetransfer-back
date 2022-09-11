package storage

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "mime/multipart"
)

type Form struct {
    Files []*multipart.FileHeader `form:"files" binding:"required"`
}

func ProcessSStorageRequest (c *gin.Context) {

    var (
        form Form
        err error
    )
    err = c.ShouldBind(&form)

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    } else {
        // for _, formFile := range form.Files {

        // Get raw file bytes - no reader method
        // openedFile, _ := formFile.Open()
        // file, _ := ioutil.ReadAll(openedFile)

        // Upload to disk
        // `formFile` has io.reader method
        // c.SaveUploadedFile(formFile, path)

        // }
        c.String(http.StatusOK, "Files uploaded")
    }
}
