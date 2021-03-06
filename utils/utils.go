package utils

import (
    "cmdb-app-mysql/models"
    "fmt"
    "strings"
    "time"

    "github.com/golang-jwt/jwt"
    "golang.org/x/crypto/bcrypt"
)

/*
明文加密
*/
func HashPassword(password string) string {
    bytePassword := []byte(password)
    passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
    return string(passwordHash)
}

/*
验证密文
*/
func VerifyPassword(passwordHash string, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}

/*
生成token
*/
const JWT_SECRET_KEY = "$2a$10$Hxx27PDN.Eu3nWhB07Dk3eHbLqT5/0zCKR/h3FtYpne6KA9Qhj0uO"
const JWT_TOKEN_EXP = 900           // 15分钟
const JWT_REFRESH_TOKEN_EXP = 86400 // 24小时

func CreateToken(id string) (*models.Login, error) {
    token_exp := time.Now().Add(time.Second * JWT_TOKEN_EXP).Unix()
    refresh_token_exp := time.Now().Add(time.Second * JWT_REFRESH_TOKEN_EXP).Unix()

    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["id"] = id
    claims["exp"] = token_exp
    tokenString, err := token.SignedString([]byte(JWT_SECRET_KEY))
    if err != nil {
        return nil, err
    }

    refreshToken := jwt.New(jwt.SigningMethodHS256)
    refreshClaims := refreshToken.Claims.(jwt.MapClaims)
    refreshClaims["id"] = id
    refreshClaims["exp"] = refresh_token_exp
    refreshTokenString, err := refreshToken.SignedString([]byte(JWT_SECRET_KEY))
    if err != nil {
        return nil, err
    }

    login := &models.Login{
        Header:               "Authorization",
        Type:                 "Bearer ",
        Token:                tokenString,
        RefreshToken:         refreshTokenString,
        TokenExpireAt:        token_exp,
        RefreshTokenExpireAT: refresh_token_exp,
    }

    return login, nil
}

/*
从header获得token
*/
func GetTokenFromHeader(header string) string {
    if header == "" {
        return ""
    }

    token := strings.Split(header, " ")
    if len(token) != 2 || token[0] != "Bearer" {
        return ""
    }

    return token[1]
}

/*
解析token，获得payload
*/
func GetPayloadFromToken(tokenString string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(JWT_SECRET_KEY), nil
    })

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, nil
    } else {
        return nil, err
    }
}

/*
构建菜单树
*/
func BuildMenuTree(items []*models.MenuTree, parentID string) []*models.MenuTree {
    tree := make([]*models.MenuTree, 0)

    for _, item := range items {
        if item.ParentID == parentID {
            item.Children = BuildMenuTree(items, item.ID)
            tree = append(tree, item)
        }
    }

    return tree
}

/*
构建权限树
*/
func BuildPremissionTree(items []*models.PermissionTree, parentID string) []*models.PermissionTree {
    tree := make([]*models.PermissionTree, 0)

    for _, item := range items {
        if item.ParentID == parentID {
            item.Children = BuildPremissionTree(items, item.ID)
            tree = append(tree, item)
        }
    }

    return tree
}

/*
构建部门树
*/
func BuildDepartmentTree(items []*models.DepartmentTree, parentID string) []*models.DepartmentTree {
    tree := make([]*models.DepartmentTree, 0)

    for _, item := range items {
        if item.ParentID == parentID {
            item.Children = BuildDepartmentTree(items, item.ID)
            tree = append(tree, item)
        }
    }

    return tree
}
