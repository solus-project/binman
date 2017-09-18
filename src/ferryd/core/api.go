//
// Copyright © 2017 Solus Project
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

package core

import (
	"fmt"
	"libeopkg"
	"path/filepath"
)

// This file provides the public API functions which are used by ferryd
// and exposed through it's handlers.

// CreateRepo will request the creation of a new repository
func (m *Manager) CreateRepo(id string) error {
	if _, err := m.repo.CreateRepo(m.db, id); err != nil {
		return err
	}
	// Index the newly created repo
	return m.Index(id)
}

// GetRepo will grab the repository if it exists
// Note that this is a read only operation
func (m *Manager) GetRepo(id string) (*Repository, error) {
	return m.repo.GetRepo(m.db, id)
}

// AddPackages will attempt to add the named packages to the repository
func (m *Manager) AddPackages(repoID string, packages []string) error {
	repo, err := m.GetRepo(repoID)
	if err != nil {
		return err
	}

	return fmt.Errorf("not yet ported")

	//if err := repo.AddPackages(m.db, m.pool, packages); err != nil {
	//	return err
	//}

	// return m.Index(repoID)
}

// Index will cause the repository's index to be reconstructed
func (m *Manager) Index(repoID string) error {
	repo, err := m.GetRepo(repoID)
	if err != nil {
		return err
	}

	return repo.Index(m.db, m.pool)
}

// GetPackageNames will attempt to load all package names for the given
// repository.
func (m *Manager) GetPackageNames(repoID string) ([]string, error) {
	repo, err := m.GetRepo(repoID)
	if err != nil {
		return nil, err
	}

	return repo.GetPackageNames(m.db)
}

// GetPackages will return a set of packages for the package name within the
// specified repository
func (m *Manager) GetPackages(repoID, pkgName string) ([]*libeopkg.MetaPackage, error) {
	repo, err := m.GetRepo(repoID)
	if err != nil {
		return nil, err
	}

	return repo.GetPackages(m.db, m.pool, pkgName)
}

// CreateDelta will attempt to create a new delta package between the old and new IDs
func (m *Manager) CreateDelta(repoID string, oldPkg, newPkg *libeopkg.MetaPackage) (string, error) {
	repo, err := m.GetRepo(repoID)
	if err != nil {
		return "", err
	}

	return repo.CreateDelta(m.db, oldPkg, newPkg)
}

// HasDelta will query the repository to determine if it already has the
// given delta
func (m *Manager) HasDelta(repoID, pkgID, deltaPath string) (bool, error) {
	repo, err := m.GetRepo(repoID)
	if err != nil {
		return false, err
	}
	return repo.HasDelta(m.db, pkgID, deltaPath)
}

// AddDelta will attempt to include the delta package specified by deltaPath into
// the target repository
func (m *Manager) AddDelta(repoID, deltaPath string, mapping *DeltaInformation) error {
	repo, err := m.GetRepo(repoID)
	if err != nil {
		return err
	}

	return repo.AddDelta(m.db, m.pool, deltaPath, mapping)
}

// MarkDeltaFailed will permanently record the delta package as failing so we do
// not attempt to recreate it (expensive)
func (m *Manager) MarkDeltaFailed(deltaID string, delta *DeltaInformation) error {
	return m.pool.MarkDeltaFailed(m.db, deltaID, delta)
}

// GetDeltaFailed will determine via the pool transaction whether a delta has
// previously failed.
func (m *Manager) GetDeltaFailed(deltaID string) bool {
	return m.pool.GetDeltaFailed(m.db, deltaID)
}

// GetPoolEntry will return the metadata for a pool entry with the given pkg ID
func (m *Manager) GetPoolEntry(pkgID string) (*libeopkg.MetaPackage, error) {
	entry, err := m.pool.GetEntry(m.db, filepath.Base(pkgID))
	if err != nil {
		return nil, err
	}
	return entry.Meta, err
}
