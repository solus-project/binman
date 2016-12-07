//
// Copyright © 2016 Ikey Doherty <ikey@solus-project.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package manager

import (
	"bytes"
	"encoding/json"
	"github.com/boltdb/bolt"
	"path/filepath"
)

var (
	// RepoDirectory is the base directory for all repositories.
	RepoDirectory = "repo"
)

// A Repository is the base unit of storage in binman
type Repository struct {
	Name string
}

// GetDirectory will return the directory component for where this
// repository lives on disk.
func (r *Repository) GetDirectory() string {
	return filepath.Join(RepoDirectory, r.Name)
}

// CreateRepo will attempt to create a new repository
func (m *Manager) CreateRepo(name string) error {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	repo := &Repository{
		Name: name,
	}
	if err := enc.Encode(repo); err != nil {
		return err
	}

	// Encoded name
	nom := []byte(name)

	return m.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketNameRepos)
		// Check it doesn't already exist in the bucket
		if b.Get(nom) != nil {
			return ErrResourceExists
		}
		return tx.Bucket(BucketNameRepos).Put(nom, buf.Bytes())
	})
}

// ListRepos will return a list of repository names known to binman.
func (m *Manager) ListRepos() ([]string, error) {
	var repos []string
	err := m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketNameRepos)
		return b.ForEach(func(k, v []byte) error {
			repos = append(repos, string(k))
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return repos, nil
}

// RemoveRepo will remove a repository from binman. In future this will also
// have to request the pool check for all unreferenced files and delete them
// too.
func (m *Manager) RemoveRepo(name string) error {
	err := m.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(BucketNameRepos)
		nom := []byte(name)
		if b.Get(nom) == nil {
			return ErrUnknownResource
		}
		return tx.Bucket(BucketNameRepos).Delete(nom)
	})
	return err
}

// GetRepo will attempt to grab the named repo, if it exists.
func (m *Manager) GetRepo(name string) (*Repository, error) {
	repo := &Repository{}
	nom := []byte(name)

	err := m.db.View(func(tx *bolt.Tx) error {
		blob := tx.Bucket(BucketNameRepos).Get(nom)
		if blob == nil {
			return ErrUnknownResource
		}
		return json.Unmarshal(blob, repo)
	})
	if err != nil {
		return nil, err
	}
	return repo, nil
}
