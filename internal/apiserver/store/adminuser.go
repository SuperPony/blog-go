/*
 * Copyright 2021 SuperPony <superponyyy@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package store

import (
	"context"

	v1 "blog-api/internal/pkg/model/v1"
	metav1 "blog-api/pkg/meta/v1"
)

type AdminUserStore interface {
	Create(ctx context.Context, adminUserModel *v1.AdminUser, opts metav1.CreateOptions) error
}