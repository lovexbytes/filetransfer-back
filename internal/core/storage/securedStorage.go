package storage

import (
	"errors"
	"strings"
	"strconv"
	"github.com/google/uuid"
	"filetransfer-back/internal/config"
	"filetransfer-back/internal/core/archive"
	"filetransfer-back/internal/core/encryption"
)

var key = []byte("jo8aegPKlvAAEa2M3l5sePv8JJMnqlTD")

type SecuredStorage struct {
  id uuid.UUID
  efiles []archive.EntryFile
  bytedata []byte
  isEncrypted bool
  isArchived bool
	exists  bool
  archiver archive.Archiver
	cipher encryption.Cipher
}

type SecuredStorageOption func (*SecuredStorage)

func NewSecuredStorage (opts ...SecuredStorageOption) (*SecuredStorage, error) {
	var nilUUID uuid.UUID

  ss := &SecuredStorage{
    id: nilUUID,
    efiles : nil,
    bytedata: nil,
  }

  for _, opt := range opts {
    opt(ss)
  }

  if !ss.exists{
		ss.id = uuid.New()
  }
  
  vErr := ss.validate()
  if (vErr != nil){
    ss = nil
  }

  return ss, vErr
}

func (ss *SecuredStorage) validate () error {
	var nilUUID uuid.UUID
  var err error = nil
  
  //id is required
  if len(uuid.UUID.String(ss.id)) == 0 {
    err = errors.New("Empty id is not allowed")
  } else {
		if ss.id == nilUUID {
			err = errors.New("nil id is not allowed")
		}
	}
	

  //files are required for the new ss
  if !ss.exists && len(ss.efiles) == 0 {
    err = errors.New("Empty efiles array is not allowed")
  }

  //archiver & cipher are required
  if ss.archiver == nil {
  	err = errors.New("Archiver cannot be nil")
  }

	if ss.cipher == nil {
  	err = errors.New("Cipher cannot be nil")
  }

  //bytedata for existing ss must be present on init
  if ss.exists && (ss.bytedata == nil || len(ss.bytedata) == 0) {
  	err = errors.New("Bytedata must be present for the existing SS")
  }

  return err
}

//options
func WithSSUUID (id uuid.UUID) SecuredStorageOption {
	return func (ss *SecuredStorage) {
		ss.id = id
		ss.exists = true
		ss.isArchived = true
		ss.isEncrypted = true
	}
}

func WithSSFileEntries (efiles []archive.EntryFile) SecuredStorageOption {
  return func (ss *SecuredStorage) {
    ss.efiles = efiles 
  }
}

func WithSSArchiver (archiver archive.Archiver) SecuredStorageOption {
	return func (ss *SecuredStorage) {
		ss.archiver = archiver
	}
}

func WithSSCipher (cipher encryption.Cipher) SecuredStorageOption {
	return func (ss *SecuredStorage) {
		ss.cipher = cipher
	}
}

func WithByteData (data []byte) SecuredStorageOption {
	return func (ss *SecuredStorage) {
		ss.bytedata = data
	}
}

//utils
func (ss *SecuredStorage) Path () string {
	return strings.Join([]string{ ss.id.String(), config.TarExtension(), ss.archiver.Extension() }, ".")
}

//getters
func (ss *SecuredStorage) Id () string {
  return uuid.UUID.String(ss.id)
}

func (ss *SecuredStorage) ByteData () ([]byte, error) {
	var data []byte
	var err error
	if ss.isArchived && ss.isEncrypted{
		data = ss.bytedata
	} else {
		err = errors.New("SS must be archived and encrypted first before getting bytedata")
	}
	return data, err
}

func (ss *SecuredStorage) EntryFiles () []archive.EntryFile {
  return ss.efiles
}

//methods
func (ss *SecuredStorage) Archive() error {
	var err error
	if !ss.isArchived {
		ss.bytedata, err = ss.archiver.Archive(ss.efiles)
		if err == nil {
			ss.isArchived = true
			//remove data of the initial files
			ss.efiles = nil
		}	
	}else{
		err = errors.New("SS is already archived")
	}
  return err
}

func (ss *SecuredStorage) Unarchive() error {
	var err error 
	if !ss.isEncrypted {
		if ss.isArchived {
			ss.efiles, err = ss.archiver.Unarchive(ss.bytedata)
			if err == nil {
				ss.isArchived = false
			}
		} else {
			err = errors.New("SS must be archived to use unarchive operation")
		}
	} else {
		err = errors.New("SS must be unencrypted before unarchive operation")
	}
	return err
}

func (ss *SecuredStorage) Encrypt () error {
	var err error
	if ss.isArchived {
		if !ss.isEncrypted {
			/**
				@TODO
				GET KEY FROM THE SECURED PLACE
			*/
			ss.bytedata, err = ss.cipher.Encrypt(key, ss.bytedata)

			if err == nil {
				ss.isEncrypted = true
			}
		}else{
			err = errors.New("SS is alredy encrypted")
		}

	} else {
		err = errors.New("SS must be archived before encryption")
	} 
	return err
}

func (ss *SecuredStorage) Decrypt () error {
	var err error

	if ss.isArchived && ss.isEncrypted {
		ss.bytedata, err = ss.cipher.Decrypt(key, ss.bytedata)
		if err == nil {
			ss.isEncrypted = false
		}
	} else {
		err = errors.New("SS must be compressed AND encrypted before decryption. SS isArchived->" + strconv.FormatBool(ss.isArchived) + ", SS isEncrypted->" + strconv.FormatBool(ss.isEncrypted))
	}
	return err
}

func (ss *SecuredStorage) Unpack () error {
	var err error
	if ss.bytedata != nil && len(ss.bytedata) > 0 {
		err = ss.Decrypt()
		if err == nil {
			err = ss.Unarchive()
		}
	} else {
		err = errors.New("bytedata empty during unpacking init")
	}
	return err
}
