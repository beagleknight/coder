package dbcrypt

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/coder/coder/v2/coderd/database"
	"github.com/coder/coder/v2/coderd/database/dbgen"
	"github.com/coder/coder/v2/coderd/database/dbtestutil"
)

func TestUserLinks(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("InsertUserLink", func(t *testing.T) {
		t.Parallel()
		db, crypt, ciphers := setup(t)
		user := dbgen.User(t, crypt, database.User{})
		link := dbgen.UserLink(t, crypt, database.UserLink{
			UserID:            user.ID,
			OAuthAccessToken:  "access",
			OAuthRefreshToken: "refresh",
		})
		require.Equal(t, "access", link.OAuthAccessToken)
		require.Equal(t, "refresh", link.OAuthRefreshToken)
		require.Equal(t, ciphers[0].HexDigest(), link.OAuthAccessTokenKeyID.String)
		require.Equal(t, ciphers[0].HexDigest(), link.OAuthRefreshTokenKeyID.String)

		rawLink, err := db.GetUserLinkByLinkedID(ctx, link.LinkedID)
		require.NoError(t, err)
		requireEncryptedEquals(t, ciphers[0], rawLink.OAuthAccessToken, "access")
		requireEncryptedEquals(t, ciphers[0], rawLink.OAuthRefreshToken, "refresh")
	})

	t.Run("UpdateUserLink", func(t *testing.T) {
		t.Parallel()
		db, crypt, ciphers := setup(t)
		user := dbgen.User(t, crypt, database.User{})
		link := dbgen.UserLink(t, crypt, database.UserLink{
			UserID: user.ID,
		})

		updated, err := crypt.UpdateUserLink(ctx, database.UpdateUserLinkParams{
			OAuthAccessToken:  "access",
			OAuthRefreshToken: "refresh",
			UserID:            link.UserID,
			LoginType:         link.LoginType,
		})
		require.NoError(t, err)
		require.Equal(t, "access", updated.OAuthAccessToken)
		require.Equal(t, "refresh", updated.OAuthRefreshToken)
		require.Equal(t, ciphers[0].HexDigest(), link.OAuthAccessTokenKeyID.String)
		require.Equal(t, ciphers[0].HexDigest(), link.OAuthRefreshTokenKeyID.String)

		rawLink, err := db.GetUserLinkByLinkedID(ctx, link.LinkedID)
		require.NoError(t, err)
		requireEncryptedEquals(t, ciphers[0], rawLink.OAuthAccessToken, "access")
		requireEncryptedEquals(t, ciphers[0], rawLink.OAuthRefreshToken, "refresh")
	})

	t.Run("GetUserLinkByLinkedID", func(t *testing.T) {
		t.Parallel()
		t.Run("OK", func(t *testing.T) {
			db, crypt, ciphers := setup(t)
			user := dbgen.User(t, crypt, database.User{})
			link := dbgen.UserLink(t, crypt, database.UserLink{
				UserID:            user.ID,
				OAuthAccessToken:  "access",
				OAuthRefreshToken: "refresh",
			})

			link, err := crypt.GetUserLinkByLinkedID(ctx, link.LinkedID)
			require.NoError(t, err)
			require.Equal(t, "access", link.OAuthAccessToken)
			require.Equal(t, "refresh", link.OAuthRefreshToken)
			require.Equal(t, ciphers[0].HexDigest(), link.OAuthAccessTokenKeyID.String)
			require.Equal(t, ciphers[0].HexDigest(), link.OAuthRefreshTokenKeyID.String)

			rawLink, err := db.GetUserLinkByLinkedID(ctx, link.LinkedID)
			require.NoError(t, err)
			requireEncryptedEquals(t, ciphers[0], rawLink.OAuthAccessToken, "access")
			requireEncryptedEquals(t, ciphers[0], rawLink.OAuthRefreshToken, "refresh")
		})

		t.Run("DecryptErr", func(t *testing.T) {
			t.Parallel()
			db, crypt, ciphers := setup(t)
			user := dbgen.User(t, db, database.User{})
			link := dbgen.UserLink(t, db, database.UserLink{
				UserID:                 user.ID,
				OAuthAccessToken:       fakeBase64RandomData(t, 32),
				OAuthRefreshToken:      fakeBase64RandomData(t, 32),
				OAuthAccessTokenKeyID:  sql.NullString{String: ciphers[0].HexDigest(), Valid: true},
				OAuthRefreshTokenKeyID: sql.NullString{String: ciphers[0].HexDigest(), Valid: true},
			})

			_, err := crypt.GetUserLinkByLinkedID(ctx, link.LinkedID)
			require.Error(t, err, "expected an error")
			var derr *DecryptFailedError
			require.ErrorAs(t, err, &derr, "expected a decrypt error")
		})
	})

	t.Run("GetUserLinksByUserID", func(t *testing.T) {
		t.Parallel()

		t.Run("OK", func(t *testing.T) {
			t.Parallel()
			db, crypt, ciphers := setup(t)
			user := dbgen.User(t, crypt, database.User{})
			link := dbgen.UserLink(t, crypt, database.UserLink{
				UserID:            user.ID,
				OAuthAccessToken:  "access",
				OAuthRefreshToken: "refresh",
			})
			links, err := crypt.GetUserLinksByUserID(ctx, link.UserID)
			require.NoError(t, err)
			require.Len(t, links, 1)
			require.Equal(t, "access", links[0].OAuthAccessToken)
			require.Equal(t, "refresh", links[0].OAuthRefreshToken)
			require.Equal(t, ciphers[0].HexDigest(), links[0].OAuthAccessTokenKeyID.String)
			require.Equal(t, ciphers[0].HexDigest(), links[0].OAuthRefreshTokenKeyID.String)

			rawLinks, err := db.GetUserLinksByUserID(ctx, link.UserID)
			require.NoError(t, err)
			require.Len(t, rawLinks, 1)
			requireEncryptedEquals(t, ciphers[0], rawLinks[0].OAuthAccessToken, "access")
			requireEncryptedEquals(t, ciphers[0], rawLinks[0].OAuthRefreshToken, "refresh")
		})

		t.Run("Empty", func(t *testing.T) {
			t.Parallel()
			_, crypt, _ := setup(t)
			user := dbgen.User(t, crypt, database.User{})
			links, err := crypt.GetUserLinksByUserID(ctx, user.ID)
			require.NoError(t, err)
			require.Empty(t, links)
		})

		t.Run("DecryptErr", func(t *testing.T) {
			t.Parallel()
			db, crypt, ciphers := setup(t)
			user := dbgen.User(t, db, database.User{})
			_ = dbgen.UserLink(t, db, database.UserLink{
				UserID:                 user.ID,
				OAuthAccessToken:       fakeBase64RandomData(t, 32),
				OAuthRefreshToken:      fakeBase64RandomData(t, 32),
				OAuthAccessTokenKeyID:  sql.NullString{String: ciphers[0].HexDigest(), Valid: true},
				OAuthRefreshTokenKeyID: sql.NullString{String: ciphers[0].HexDigest(), Valid: true},
			})
			_, err := crypt.GetUserLinksByUserID(ctx, user.ID)
			require.Error(t, err, "expected an error")
			var derr *DecryptFailedError
			require.ErrorAs(t, err, &derr, "expected a decrypt error")
		})
	})

	t.Run("GetUserLinkByUserIDLoginType", func(t *testing.T) {
		t.Parallel()
		t.Run("OK", func(t *testing.T) {
			db, crypt, ciphers := setup(t)
			user := dbgen.User(t, crypt, database.User{})
			link := dbgen.UserLink(t, crypt, database.UserLink{
				UserID:            user.ID,
				OAuthAccessToken:  "access",
				OAuthRefreshToken: "refresh",
			})

			link, err := crypt.GetUserLinkByUserIDLoginType(ctx, database.GetUserLinkByUserIDLoginTypeParams{
				UserID:    link.UserID,
				LoginType: link.LoginType,
			})
			require.NoError(t, err)
			require.Equal(t, "access", link.OAuthAccessToken)
			require.Equal(t, "refresh", link.OAuthRefreshToken)
			require.Equal(t, ciphers[0].HexDigest(), link.OAuthAccessTokenKeyID.String)
			require.Equal(t, ciphers[0].HexDigest(), link.OAuthRefreshTokenKeyID.String)

			rawLink, err := db.GetUserLinkByUserIDLoginType(ctx, database.GetUserLinkByUserIDLoginTypeParams{
				UserID:    link.UserID,
				LoginType: link.LoginType,
			})
			require.NoError(t, err)
			requireEncryptedEquals(t, ciphers[0], rawLink.OAuthAccessToken, "access")
			requireEncryptedEquals(t, ciphers[0], rawLink.OAuthRefreshToken, "refresh")
		})

		t.Run("DecryptErr", func(t *testing.T) {
			db, crypt, ciphers := setup(t)
			user := dbgen.User(t, db, database.User{})
			link := dbgen.UserLink(t, db, database.UserLink{
				UserID:                 user.ID,
				OAuthAccessToken:       fakeBase64RandomData(t, 32),
				OAuthRefreshToken:      fakeBase64RandomData(t, 32),
				OAuthAccessTokenKeyID:  sql.NullString{String: ciphers[0].HexDigest(), Valid: true},
				OAuthRefreshTokenKeyID: sql.NullString{String: ciphers[0].HexDigest(), Valid: true},
			})

			_, err := crypt.GetUserLinkByUserIDLoginType(ctx, database.GetUserLinkByUserIDLoginTypeParams{
				UserID:    link.UserID,
				LoginType: link.LoginType,
			})
			require.Error(t, err, "expected an error")
			var derr *DecryptFailedError
			require.ErrorAs(t, err, &derr, "expected a decrypt error")
		})
	})
}

func TestGitAuthLinks(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("InsertGitAuthLink", func(t *testing.T) {
		t.Parallel()
		db, crypt, ciphers := setup(t)
		link := dbgen.GitAuthLink(t, crypt, database.GitAuthLink{
			OAuthAccessToken:  "access",
			OAuthRefreshToken: "refresh",
		})
		require.Equal(t, "access", link.OAuthAccessToken)
		require.Equal(t, "refresh", link.OAuthRefreshToken)

		link, err := db.GetGitAuthLink(ctx, database.GetGitAuthLinkParams{
			ProviderID: link.ProviderID,
			UserID:     link.UserID,
		})
		require.NoError(t, err)
		requireEncryptedEquals(t, ciphers[0], link.OAuthAccessToken, "access")
		requireEncryptedEquals(t, ciphers[0], link.OAuthRefreshToken, "refresh")
	})

	t.Run("UpdateGitAuthLink", func(t *testing.T) {
		t.Parallel()
		db, crypt, ciphers := setup(t)
		link := dbgen.GitAuthLink(t, crypt, database.GitAuthLink{})
		updated, err := crypt.UpdateGitAuthLink(ctx, database.UpdateGitAuthLinkParams{
			ProviderID:        link.ProviderID,
			UserID:            link.UserID,
			OAuthAccessToken:  "access",
			OAuthRefreshToken: "refresh",
		})
		require.NoError(t, err)
		require.Equal(t, "access", updated.OAuthAccessToken)
		require.Equal(t, "refresh", updated.OAuthRefreshToken)

		link, err = db.GetGitAuthLink(ctx, database.GetGitAuthLinkParams{
			ProviderID: link.ProviderID,
			UserID:     link.UserID,
		})
		require.NoError(t, err)
		requireEncryptedEquals(t, ciphers[0], link.OAuthAccessToken, "access")
		requireEncryptedEquals(t, ciphers[0], link.OAuthRefreshToken, "refresh")
	})

	t.Run("GetGitAuthLink", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			t.Parallel()
			db, crypt, ciphers := setup(t)
			link := dbgen.GitAuthLink(t, crypt, database.GitAuthLink{
				OAuthAccessToken:  "access",
				OAuthRefreshToken: "refresh",
			})
			link, err := db.GetGitAuthLink(ctx, database.GetGitAuthLinkParams{
				UserID:     link.UserID,
				ProviderID: link.ProviderID,
			})
			require.NoError(t, err)
			requireEncryptedEquals(t, ciphers[0], link.OAuthAccessToken, "access")
			requireEncryptedEquals(t, ciphers[0], link.OAuthRefreshToken, "refresh")
		})
		t.Run("DecryptErr", func(t *testing.T) {
			t.Parallel()
			db, crypt, ciphers := setup(t)
			link := dbgen.GitAuthLink(t, db, database.GitAuthLink{
				OAuthAccessToken:       fakeBase64RandomData(t, 32),
				OAuthRefreshToken:      fakeBase64RandomData(t, 32),
				OAuthAccessTokenKeyID:  sql.NullString{String: ciphers[0].HexDigest(), Valid: true},
				OAuthRefreshTokenKeyID: sql.NullString{String: ciphers[0].HexDigest(), Valid: true},
			})

			_, err := crypt.GetGitAuthLink(ctx, database.GetGitAuthLinkParams{
				UserID:     link.UserID,
				ProviderID: link.ProviderID,
			})
			require.Error(t, err, "expected an error")
			var derr *DecryptFailedError
			require.ErrorAs(t, err, &derr, "expected a decrypt error")
		})
	})

	t.Run("GetGitAuthLinksByUserID", func(t *testing.T) {
		t.Parallel()

		t.Run("OK", func(t *testing.T) {
			t.Parallel()
			db, crypt, ciphers := setup(t)
			user := dbgen.User(t, crypt, database.User{})
			link := dbgen.GitAuthLink(t, crypt, database.GitAuthLink{
				UserID:            user.ID,
				OAuthAccessToken:  "access",
				OAuthRefreshToken: "refresh",
			})
			links, err := crypt.GetGitAuthLinksByUserID(ctx, link.UserID)
			require.NoError(t, err)
			require.Len(t, links, 1)
			require.Equal(t, "access", links[0].OAuthAccessToken)
			require.Equal(t, "refresh", links[0].OAuthRefreshToken)
			require.Equal(t, ciphers[0].HexDigest(), links[0].OAuthAccessTokenKeyID.String)
			require.Equal(t, ciphers[0].HexDigest(), links[0].OAuthRefreshTokenKeyID.String)

			rawLinks, err := db.GetGitAuthLinksByUserID(ctx, link.UserID)
			require.NoError(t, err)
			require.Len(t, rawLinks, 1)
			requireEncryptedEquals(t, ciphers[0], rawLinks[0].OAuthAccessToken, "access")
			requireEncryptedEquals(t, ciphers[0], rawLinks[0].OAuthRefreshToken, "refresh")
		})

		t.Run("DecryptErr", func(t *testing.T) {
			db, crypt, ciphers := setup(t)
			user := dbgen.User(t, db, database.User{})
			link := dbgen.GitAuthLink(t, db, database.GitAuthLink{
				UserID:                 user.ID,
				OAuthAccessToken:       fakeBase64RandomData(t, 32),
				OAuthRefreshToken:      fakeBase64RandomData(t, 32),
				OAuthAccessTokenKeyID:  sql.NullString{String: ciphers[0].HexDigest(), Valid: true},
				OAuthRefreshTokenKeyID: sql.NullString{String: ciphers[0].HexDigest(), Valid: true},
			})
			_, err := crypt.GetGitAuthLinksByUserID(ctx, link.UserID)
			require.Error(t, err, "expected an error")
			var derr *DecryptFailedError
			require.ErrorAs(t, err, &derr, "expected a decrypt error")
		})
	})
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()
		// Given: a cipher is loaded
		cipher := initCipher(t)
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		rawDB, _ := dbtestutil.NewDB(t)

		// Before: no keys should be present
		keys, err := rawDB.GetDBCryptKeys(ctx)
		require.NoError(t, err, "no error should be returned")
		require.Empty(t, keys, "no keys should be present")

		// When: we init the crypt db
		_, err = New(ctx, rawDB, cipher)
		require.NoError(t, err)

		// Then: a new key is inserted
		keys, err = rawDB.GetDBCryptKeys(ctx)
		require.NoError(t, err)
		require.Len(t, keys, 1, "one key should be present")
		require.Equal(t, cipher.HexDigest(), keys[0].ActiveKeyDigest.String, "key digest mismatch")
		require.Empty(t, keys[0].RevokedKeyDigest.String, "key should not be revoked")
		requireEncryptedEquals(t, cipher, keys[0].Test, "coder")
	})

	t.Run("MissingKey", func(t *testing.T) {
		t.Parallel()

		// Given: there exist two valid encryption keys
		cipher1 := initCipher(t)
		cipher2 := initCipher(t)
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		rawDB, _ := dbtestutil.NewDB(t)

		// Given: key 1 is already present in the database
		err := rawDB.InsertDBCryptKey(ctx, database.InsertDBCryptKeyParams{
			Number:          1,
			ActiveKeyDigest: cipher1.HexDigest(),
			Test:            fakeBase64RandomData(t, 32),
		})
		require.NoError(t, err, "no error should be returned")
		keys, err := rawDB.GetDBCryptKeys(ctx)
		require.NoError(t, err, "no error should be returned")
		require.Len(t, keys, 1, "one key should be present")

		// When: we init the crypt db with key 2
		_, err = New(ctx, rawDB, cipher2)

		// Then: we error because the key is not revoked and we don't know how to decrypt it
		require.Error(t, err)
		var derr *DecryptFailedError
		require.ErrorAs(t, err, &derr, "expected a decrypt error")

		// When: the existing key is marked as having been revoked
		err = rawDB.RevokeDBCryptKey(ctx, cipher1.HexDigest())
		require.NoError(t, err, "no error should be returned")

		// And: we init the crypt db with key 2
		_, err = New(ctx, rawDB, cipher2)

		// Then: we succeed
		require.NoError(t, err)

		// And: key 2 should now be the active key
		keys, err = rawDB.GetDBCryptKeys(ctx)
		require.NoError(t, err)
		require.Len(t, keys, 2, "two keys should be present")
		require.EqualValues(t, keys[0].Number, 1, "key number mismatch")
		require.Empty(t, keys[0].ActiveKeyDigest.String, "key should not be active")
		require.Equal(t, cipher1.HexDigest(), keys[0].RevokedKeyDigest.String, "key should be revoked")

		require.EqualValues(t, keys[1].Number, 2, "key number mismatch")
		require.Equal(t, cipher2.HexDigest(), keys[1].ActiveKeyDigest.String, "key digest mismatch")
		require.Empty(t, keys[1].RevokedKeyDigest.String, "key should not be revoked")
		requireEncryptedEquals(t, cipher2, keys[1].Test, "coder")
	})

	t.Run("NoCipher", func(t *testing.T) {
		t.Parallel()
		// Given: no cipher is loaded
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		rawDB, _ := dbtestutil.NewDB(t)

		keys, err := rawDB.GetDBCryptKeys(ctx)
		require.NoError(t, err, "no error should be returned")
		require.Empty(t, keys, "no keys should be present")

		// When: we init the crypt db with no ciphers
		cs := make([]Cipher, 0)
		_, err = New(ctx, rawDB, cs...)

		// Then: an error is returned
		require.ErrorContains(t, err, "no ciphers configured")

		// Assert invariant: no keys are inserted
		keys, err = rawDB.GetDBCryptKeys(ctx)
		require.NoError(t, err, "no error should be returned")
		require.Empty(t, keys, "no keys should be present")
	})

	t.Run("PrimaryRevoked", func(t *testing.T) {
		t.Parallel()
		// Given: a cipher is loaded
		cipher := initCipher(t)
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		rawDB, _ := dbtestutil.NewDB(t)

		// And: the cipher is revoked before we init the crypt db
		err := rawDB.InsertDBCryptKey(ctx, database.InsertDBCryptKeyParams{
			Number:          1,
			ActiveKeyDigest: cipher.HexDigest(),
			Test:            fakeBase64RandomData(t, 32),
		})
		require.NoError(t, err, "no error should be returned")
		err = rawDB.RevokeDBCryptKey(ctx, cipher.HexDigest())
		require.NoError(t, err, "no error should be returned")

		// Then: when we init the crypt db, we error because the key is revoked
		_, err = New(ctx, rawDB, cipher)
		require.Error(t, err)
		require.ErrorContains(t, err, "has been revoked")
	})
}

func requireEncryptedEquals(t *testing.T, c Cipher, value, expected string) {
	t.Helper()
	data, err := base64.StdEncoding.DecodeString(value)
	require.NoError(t, err, "invalid base64")
	got, err := c.Decrypt(data)
	require.NoError(t, err, "failed to decrypt data")
	require.Equal(t, expected, string(got), "decrypted data does not match")
}

func initCipher(t *testing.T) *aes256 {
	t.Helper()
	key := make([]byte, 32) // AES-256 key size is 32 bytes
	_, err := io.ReadFull(rand.Reader, key)
	require.NoError(t, err)
	c, err := cipherAES256(key)
	require.NoError(t, err)
	return c
}

func setup(t *testing.T) (db, cryptodb database.Store, cs []Cipher) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	rawDB, _ := dbtestutil.NewDB(t)

	cs = append(cs, initCipher(t))
	cryptDB, err := New(ctx, rawDB, cs...)
	require.NoError(t, err)

	return rawDB, cryptDB, cs
}

func fakeBase64RandomData(t *testing.T, n int) string {
	t.Helper()
	b := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, b)
	require.NoError(t, err)
	return base64.StdEncoding.EncodeToString(b)
}
