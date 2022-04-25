package services

import (
    "cmdb-app-mysql/models"
    "cmdb-app-mysql/utils"
    "time"
)

/* 用户登录 */
func (u *UserServiceImpl) LoginUser(user *models.UserLogin) (*models.Login, error) {
    var userID, hashPassword string
    var adminFlag int8

    sql := `
    select id, password, admin_flag
    from sys_user
    where
        email = ? or mobile = ? and status = 1
    `

    row := u.mysqlClient.QueryRowContext(u.ctx, sql, user.LoginID, user.LoginID)
    if err := row.Scan(&userID, &hashPassword, &adminFlag); err != nil {
        return nil, err
    }

    // 校验密码
    if err := utils.VerifyPassword(hashPassword, user.Password); err != nil {
        return nil, err
    }

    // 生成token
    login, err := utils.CreateToken(userID, adminFlag)
    if err != nil {
        return nil, err
    }

    // 缓存token
    if err := u.WriteToRedis(userID, login.Token, time.Second*utils.JWT_TOKEN_EXP); err != nil {
        return nil, err
    }

    return login, nil
}

/* 用户注销 */
func (u *UserServiceImpl) LogoutUser() error {
    return nil
}

/* 用户刷新 */
func (u *UserServiceImpl) RefreshUser() (*models.Login, error) {
    return nil, nil
}

/* 用户变更密码 */
func (u *UserServiceImpl) ChangePassword(user *models.PasswordChange) error {
    return nil
}

/* 用户重置密码 */
func (u *UserServiceImpl) ResetPassword(user *models.PasswordReset) error {
    return nil
}

/* 从redis获取token */
func (u *UserServiceImpl) ReadFromRedis(key string) (string, error) {
    cmd := u.redisClient.Get(u.ctx, key)
    if cmd.Err() != nil {
        return "", cmd.Err()
    }

    return cmd.Val(), nil
}

/* 向redis写入token */
func (u *UserServiceImpl) WriteToRedis(key string, value string, expiration time.Duration) error {
    if err := u.redisClient.Set(u.ctx, key, value, expiration).Err(); err != nil {
        return err
    }
    return nil
}

/* 从redis删除token */
func (u *UserServiceImpl) RemoveFromRedis(key string) error {
    if err := u.redisClient.Del(u.ctx, key).Err(); err != nil {
        return err
    }
    return nil
}