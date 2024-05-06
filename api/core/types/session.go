package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

const LoginUserID = "LOGIN_USER_ID"
const LoginUserCache = "LOGIN_USER_CACHE"

const UserAuthHeader = "Authorization"
const AdminAuthHeader = "Admin-Authorization"

// Session configs struct
type Session struct {
	SecretKey string
	MaxAge    int
}
