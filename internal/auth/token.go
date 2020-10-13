package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"

	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"github.com/google/uuid"
)

var (
	ErrInvalid = fmt.Errorf("%w: invalid token", ErrAuth)
)

// TODO add a timestamp
func TokenForID(tx db.Tx, id db.ID) (string, error) {
	mac, err := getMac(tx)
	if err != nil {
		return "", err
	}
	bytes, err := id.Bytes()

	if err != nil {
		return "", err
	}
	mac.Write(bytes)
	bytesMac := mac.Sum(nil)
	binaryToken := append(bytes, bytesMac...)
	token := base64.StdEncoding.EncodeToString(binaryToken)
	return token, nil
}

func UserForToken(tx db.Tx, b64Token string) (db.Record, error) {
	binaryToken, err := base64.StdEncoding.DecodeString(b64Token)
	if err != nil {
		return nil, err
	}
	uuidBytes := binaryToken[:16]
	providedMacBytes := binaryToken[16:]

	mac, err := getMac(tx)
	mac.Write(uuidBytes)
	computedMacBytes := mac.Sum(nil)

	if !hmac.Equal(providedMacBytes, computedMacBytes) {
		return nil, ErrInvalid
	}
	id, err := uuid.FromBytes(uuidBytes)
	if err != nil {
		return nil, err
	}

	users := tx.Ref(UserModel.ID())
	user, err := tx.Query(users, db.Filter(users, db.EqID(db.ID(id)))).OneRecord()
	if err != nil {
		return nil, ErrInvalid
	}
	return user, nil
}

var initializeAuthKey = func(event lib.DatabaseReady) {
	appDB := event.Db
	tx := appDB.NewTxWithContext(noAuthContext)

	keys := tx.Ref(AuthKeyModel.ID())
	rec, err := tx.Query(keys, db.Filter(keys, db.Eq("active", true))).OneRecord()

	if errors.Is(db.ErrNotFound, err) {
		rec, err = createAuthKey()
		rwtx := appDB.NewRWTx()
		rwtx.Insert(rec)
		rwtx.Commit()
	}
}

func getMac(tx db.Tx) (hash.Hash, error) {
	keys := tx.Ref(AuthKeyModel.ID())
	rec, err := tx.Query(keys, db.Filter(keys, db.Eq("active", true))).OneRecord()
	if err != nil {
		return nil, err
	}

	b64KeyIf, err := rec.Get("key")
	if err != nil {
		return nil, err
	}
	b64Key := b64KeyIf.(string)
	key, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		return nil, err
	}
	mac := hmac.New(sha256.New, key)
	return mac, nil
}

func createAuthKey() (db.Record, error) {
	akStore := db.RecordForModel(AuthKeyModel)
	// 128-bit key
	// https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html#Session_ID_Length
	c := 16
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	b64Key := base64.StdEncoding.EncodeToString(b)

	akStore.Set("id", uuid.New())
	akStore.Set("active", true)
	akStore.Set("key", b64Key)

	return akStore, nil
}
