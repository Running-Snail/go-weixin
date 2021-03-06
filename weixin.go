package weixin

import (
    "fmt"
    "log"
    "encoding/json"
    "net/http"
    "net/url"
    "io/ioutil"
    "crypto/sha1"
)

type Weixin struct {
    AppId string
    AppSecret string
}

type WebAccessTokenResponse struct {
    AccessToken string `json:"access_token"`
    ExpiresIn int64 `json:"expires_in"`
    RefreshToken string `json:"refresh_token"`
    OpenId string `json:"openid"`
    Scope string `json:"scope"`
    UnionId string `json:"unionid"`
    Errcode int64 `json:"errcode"`
    Errmsg string `json:"errmsg"`
}

func (resp *WebAccessTokenResponse) Ok() bool {
    if resp.Errcode == 0 {
        return true
    }
    return false
}

type AccessTokenResponse struct {
    AccessToken string `json:"access_token"`
    ExpiresIn int64 `json:"expires_in"`
    Errcode int64 `json:"errcode"`
    Errmsg string `json:"errmsg"`
}

func (resp *AccessTokenResponse) Ok() bool {
    if resp.Errcode == 0 {
        return true
    }
    return false
}

type UserInfoResponse struct {
    OpenId string `json:"openid"`
    Nickname string `json:"nickname"`
    Sex int64 `json:"sex"`
    Province string `json:"province"`
    City string `json:"city"`
    Country string `json:"country"`
    Headimgurl string `json:"headimgurl"`
    Privilege []string `json:"privilege"`
    UnionId string `json:"unionid"`
    Errcode int64 `json:"errcode"`
    Errmsg string `json:"errmsg"`
}

func (resp *UserInfoResponse) Ok() bool {
    if resp.Errcode == 0 {
        return true
    }
    return false
}

type JSSDKTicketResponse struct {
    Ticket string `json:"ticket"`
    ExpiresIn int `json:"expires_in"`
    Errcode int64 `json:"errcode"`
    Errmsg string `json:"errmsg"`
}

func (resp *JSSDKTicketResponse) Ok() bool {
    if resp.Errcode == 0 {
        return true
    }
    return false
}

const (
    webAuthRedirectURL = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect"
    getAccessToken = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
    getWebAccessToken = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
    getUserInfo = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
    getJSSDKTicket = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
)

func New(appId string, appSecret string) Weixin {
    return Weixin{
        AppId: appId,
        AppSecret: appSecret,
    }
}

func (wx *Weixin) WebAuthRedirectURL(redirectURI string, scope string, state string) string {
    redirectURIEscaped := url.QueryEscape(redirectURI)
    return fmt.Sprintf(webAuthRedirectURL, wx.AppId, redirectURIEscaped, scope, state)
}

func requestGet(url string) ([]byte, error) {
    log.Println("get response on url: ", url)
    resp, err := http.Get(url)
    if err != nil {
        log.Println("failed to get url")
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Println("failed to read response body")
        return nil, err
    }
    log.Println("response body is ", string(body))

    return body, nil
}

func (wx *Weixin) GetJSSDKTicket(accessToken string) (JSSDKTicketResponse, error) {
    var response JSSDKTicketResponse
    url := fmt.Sprintf(getJSSDKTicket, accessToken)
    log.Println("get access token request url: ", url)
    body, err := requestGet(url)
    if err != nil {
        return response, err
    }

    err = json.Unmarshal(body, &response)
    if err != nil {
        log.Println("failed to parse body to json")
        return response, err
    }
    log.Printf("body json response is %v\n", response)
    return response, nil
}

func (wx *Weixin) GetAccessToken() (AccessTokenResponse, error) {
    var response AccessTokenResponse
    url := fmt.Sprintf(getAccessToken, wx.AppId, wx.AppSecret)
    log.Println("get access token request url: ", url)
    body, err := requestGet(url)
    if err != nil {
        return response, err
    }

    err = json.Unmarshal(body, &response)
    if err != nil {
        log.Println("failed to parse body to json")
        return response, err
    }
    log.Printf("body json response is %v\n", response)
    return response, nil
}

func (wx *Weixin) GetWebAccessToken(code string) (WebAccessTokenResponse, error) {
    var response WebAccessTokenResponse
    url := fmt.Sprintf(getWebAccessToken, wx.AppId, wx.AppSecret, code)
    log.Println("get web access token request url: ", url)
    body, err := requestGet(url)
    if err != nil {
        return response, err
    }
    
    err = json.Unmarshal(body, &response)
    if err != nil {
        log.Println("failed to parse body to json")
        return response, err
    }
    log.Printf("body json response is %v\n", response)
    return response, nil
}

func (wx *Weixin) GetUserInfo(accessToken string, openId string) (UserInfoResponse, error) {
    var response UserInfoResponse
    url := fmt.Sprintf(getUserInfo, accessToken, openId)
    log.Println("get user info request url: ", url)
    body, err := requestGet(url)
    if err != nil {
        return response, err
    }

    err = json.Unmarshal(body, &response)
    if err != nil {
        log.Println("failed to parse body to json")
        return response, err
    }
    log.Printf("body json response is %v\n", response)
    return response, nil
}

func (wx *Weixin) JSSDKSignature(jssdkTicket string, noncestr string, timestamp int64, url string) string {
    string1 := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", jssdkTicket, noncestr, timestamp, url)
    return fmt.Sprintf("%x", sha1.Sum([]byte(string1)))
}
