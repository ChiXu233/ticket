package service

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
)

func TestRefreshToken(t *testing.T) {
	var wg sync.WaitGroup
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://127.0.0.1:8401/api/fresh_token", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 创建一个cookie示例，这里假设名为"myCookie"，值为"cookieValue"，实际根据你真实的cookie情况调整
	cookie := &http.Cookie{
		Name:  "x-token",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVVUlEIjoiMWE4OTJhNGYtNThjNi00MDNhLTk2MTYtZDVlODgyNWIwODRhIiwiSUQiOjEsIlVzZXJJRCI6MSwiVXNlcm5hbWUiOiJhZG1pbiIsIlJvbGVOYW1lIjoic3VwZXJfYWRtaW4iLCJBdXRob3JpdHlJZCI6MCwiQnVmZmVyVGltZSI6ODY0MDAsImlzcyI6ImFkbWluIiwiZXhwIjoxNzM1MTEwODYwLCJuYmYiOjE3MzUxMTA3NDAsImlhdCI6MTczNTExMDc0MH0.gG3Bj2wbLE-c8AGOgFdmQ6YIEpXT0zAiCSnZPnEWtAM",
	}
	req.AddCookie(cookie)

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()
			defer wg.Done()
			fmt.Println(resp.Body)
		}()
	}
	wg.Wait()
	return
	// 后续处理响应获取token等逻辑，这里简化返回空字符串示例，实际按真实逻辑补充

}
