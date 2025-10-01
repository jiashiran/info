package main

import (
	"bufio"
	"fmt"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/ssh"
	sshd "info/pkg/ssh"
	"info/pkg/xueqiu"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	serverMap map[string]*ssh.Client

	//ips = []string{"39.105.202.8","39.107.243.195","39.96.50.197","123.56.17.69","47.95.241.44","47.94.130.119","47.93.218.38","101.200.52.2"}
)

func main() {
	//jijin.Run()
	xueqiu.Run()
	//areaCode.Get()
	//areaCode.Detail("france")
	//xueqiu.RunX()
}

func rep(s string) string {
	if strings.Index(s, ` `) > 0 {
		s = strings.ReplaceAll(strings.Split(s, ` `)[0], " ", "")
	}
	if strings.Index(s, `or`) > 0 {
		s = strings.ReplaceAll(strings.Split(s, `or`)[0], " ", "")
	}
	if strings.Index(s, `+`) > 0 {
		s = strings.ReplaceAll(strings.Split(s, `+`)[0], " ", "")
	}
	if strings.Index(s, `&lt;`) > 0 {
		s = strings.ReplaceAll(strings.Split(s, `&lt;`)[0], " ", "")
	}
	if strings.Index(s, `digits`) > 0 {
		s = strings.ReplaceAll(strings.Split(s, `digits`)[0], " ", "")
	}
	if strings.Index(s, `XX`) > 0 {
		s = strings.ReplaceAll(strings.Split(s, `XX`)[0], " ", "")
	}
	if strings.Index(s, `)`) > 0 {
		s = strings.ReplaceAll(strings.Split(s, `)`)[0], " ", "")
	}
	return s
}

// email verify
func VerifyEmailFormat(email string) string {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	//pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	es := reg.FindAllString(email, -1)
	/*if len(es) > 0{
		log.Println(es)
	} else {
		//log.Println(email)
	}*/

	return fmt.Sprint(es)
}

// mobile verify
func VerifyMobileFormat(mobileNum string) string {
	//regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	regular := `(\(\d{3,4}\)|\d{3,4}-|\s)?\d{7,14}` //\d+`
	reg := regexp.MustCompile(regular)
	ss := reg.FindAllString(mobileNum, -1)
	if len(ss) > 0 {
		//log.Println(ss)
	} else {
		//log.Println(mobileNum)
	}
	return fmt.Sprint(ss)
}

// 获得百度token
func getaccess_token() {
	value := make(url.Values)
	resp, _ := http.PostForm("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=loC7VdqRiDBcDtdOxsARjeZc&client_secret=l0Z5yry7YLbGnDIIqQKSKlvIKGOGY6gs&", value)
	bs, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(bs))
	// "access_token":"24.cb9afce813b0dc1449529649bfc2a3dc.2592000.1570863615.282335-17235293"
}

func post() {
	//face++
	/*value := make(url.Values)
	value.Add("api_key","-QZz6z1pR_D5_X8KQV3mFwA02PS-Z8AU")
	value.Add("api_secret","MaFmJt0LIon69262LaEsbkD28q6zVjOn")
	value.Add("image_base64","/9j/2wCEAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDIBCQkJDAsMGA0NGDIhHCEyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMv/AABEIAEACWAMBIgACEQEDEQH/xAGiAAABBQEBAQEBAQAAAAAAAAAAAQIDBAUGBwgJCgsQAAIBAwMCBAMFBQQEAAABfQECAwAEEQUSITFBBhNRYQcicRQygZGhCCNCscEVUtHwJDNicoIJChYXGBkaJSYnKCkqNDU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6g4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2drh4uPk5ebn6Onq8fLz9PX29/j5+gEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoLEQACAQIEBAMEBwUEBAABAncAAQIDEQQFITEGEkFRB2FxEyIygQgUQpGhscEJIzNS8BVictEKFiQ04SXxFxgZGiYnKCkqNTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqCg4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2dri4+Tl5ufo6ery8/T19vf4+fr/2gAMAwEAAhEDEQA/APf6PrUc7slu7qMsBkCvJ9d8Y+LItRMFtZN5OcZC12YPBVMVJqDSt3diZSUT1WS6tov9ZKi/U1VTW9KebylvYTJ/dDc18/XPinXNR8Sx6bLvRnOCK9D0L4ZNDqSancXr7uuzPevTr5PSwsE8RUs2rq2pCqOWyPTgcgEYINLmuD1/xhd+HLxYZIs2443kV0uka3a65pZntJAzFDkDsa8upg6tOmqrXuvqWpJuxcXUbF5jEtzGZB1UHmrYORmvB9HN0PH9wHnYjeflz717pD/x7p9BW2YYFYRxSle6uKEuYkopNw9aQsMda86xY4kAZPSqS6pp7XP2dbqIzf3N3NT3LD7NJz/Ca+dtJkuB8YSGuGKeZ03cYr1Muy5YuNSTlblV/Uic+Wx9H0U3cM9aXcPWvLsWLVaa+tIJRHNOiOeik81PuGeteW/FLS9RMkeqafI/7obiq12YHDRxFZUpS5b9SZOyueoySpFEZJGVYwMlj0xWSfE+ghiDqdtkdt9cN4O8WzeMdKk0m8zFJs8st0NZ958DbYiaZNVlDHLYJNd0Muw9GpKljajhJbWV7kubavFHpZ8SaGF3HUrbHrvpv/CUaB/0FLb/AL7r5x0rwTNqHjL+wn1BxGGxu3GvRv8AhQ9mM41Wb9a7cTlWXYZqNWu7tX26Eqc5bI9SttY0y7bbbXkMjeitmrxYBdxPHXNfN3gizuNH+JUul/ankjiYry3Wvo5o99uY8/eXGa83NcvhgqsYxlzJq9y4T5kZs3iLRIJCk2o26OOoLUz/AISnQP8AoKW3/fdcNq/wZtdVvnum1KVC5zgZqh/woez/AOgrN+tdEMLlTinKvJP/AAicp9j0j/hKdA/6Clt/33Tk8SaHIQE1K2J9mrzX/hQ9nn/kKzfrXHeMPhy/hWSOS31KRhkHljXRQy3LMRP2dOu7vyJc5pXaPo6GaKeMPE4dD0Ip+fSuV+Hkjv4Tt977nHUmofGvji18LWYEbq90WxsrxvqVSWJeHp6u9jTmVrs7H61Xury1s03XMyRL6scVkeFtan1zTIruVNhYZxiuL+MnnSWdrBFKY9/GQfetMNgXVxSw03buEpWjdHdnxPoIODqdsD/v1ZttY0y7bbbXkMh9FbNeL6L8FjqejpeS6pIJJFyoyapeFtEufD3i9bVrtpFD4wTXqyynAyjNUazco9LGftJaXR79cXMFtHvnlWNPVjWa3ibQVOG1O2B/365T4svIPC22KUpIwxwa4Lwd8JB4h0NdQu9SlR3JAAJNc2Ey3DTwv1nEVHFXtorlSnLmske1xeIdFmYLFqFuxPYNWmjq6BkIKnoRXzjq/gaXwrrUYh1B3AboWr3nw45bQLYu2W29azzHL6OHpxq0J8yl5WCE23Zo16hnure3GZ5UQf7RqT0PavI/i9rDQ3djZ20xV5CAdp965MBg3i66pJ2uVKXKrnrMFxDcLuhkVx6qak57Vzng7TpNO0WEzSFmkQE5rA8S/ENtI8TQ6TbKJDIQOmacMFOtXlSoa2v+Ac1ldneTXdtAcTSoh9zUkciSoGjYMp7ivO/iRa30nhb+0rVmW4ABKip/hX4gGp+HEguZP9NQ8q3WtJZf/sf1mLvZ2a7C5/e5T0CjIpBwMGsjUPE+jaY+28vYom9GrghTnUdoK78i27GvketRyzwwoWlkVVHcmue/4T3wvkD+1IK4P4mePdNNlHb6deBzIPmKGu/C5XiK9aNPkav1sRKaSuet215bXSk28ySAdSpqbcPWvJfhv4s8Paborrd6ogmc5O812g8f+Ff+grBRistrUq0qcYSaXWzCM01c6bI9aMj1rmf+E/8AC3/QWt/zqWHxv4buJAkWpwMx6AVzPB4hb039zK5l3Nm7vbSzTddTpEvqxxVA+J9BBwdTts/79eafHC5kbSLFraYiORhkqcZrM074OwXXhaPU59VkWR4DKeTgHHFevh8rwv1WGIxFVx5m0klczc5c1kj2eDV9NuATDdwuB/daoZvEOiwNibULdD6Fq8Q+F2k3tw2rI87lIQQpJJzisW20Wx1fxBcxavqzW6q5HLYrp/sGgq1SEqrtG2yu9Re1dk7H0J/wlXh7/oK2v/fdKvijQGOF1S2J/wB+vIz8OvA20H/hJef+ulZXiXwR4V0zRJLnTPEBluk+6u/rU08pwNSShGpO7/ug6kl0PoO2u7a7TfbTJIvqpzU9eYfBwzHQB5spkI6EmvT+9eHjsMsNiJUU72NIu6uFV5721tmCzzojHoGNLc3MdpbSTzMFRQTk14Pe61eeNPGQgsJX8qB+dv1rfL8veLcm3aMVqxTnynvwZSoYHKkZBqu19ZrMImnjEh6Lnmi0iKafDEx+ZYwDXiHj+fU9E8fWtxHI/wBk3gt6UZfgVjKrpKVnZtedgnLlVz3iiqOlahFqmnQ3MLBlZRnHrV37wrz5RcW4vdFkVzcwWkRkuJVjT+8x4rMPijQAcHVLb/vuneItDGv6Y1m0pjDfxCvNX+BFszlv7Wl5Oe9ejgqGBqQbxFVxfkrkSclsj0j/AISnw/8A9BS2/wC+6P8AhKfD/wD0FLb/AL7rzX/hQ1r/ANBaX9aP+FDWv/QWl/Wuz6plP/QRL/wEnmqdj0r/AISnw/8A9BS2/wC+6P8AhKNA/wCgpbf9915r/wAKGtf+gtL+tH/Ch7YDjVpc/U0fVMp/6CJf+AhzVOx6na61pd5IEtr2GRz0CtmtCvnLw1o8/hv4mxWC3jSxq+OW619GHqK5M0wFPCTiqcuZSV7lQk5LUWiiivLLCiiigBkkixxNI5wqjJrlL34h+FbKZori7UOOoKCuqmiWaF4m+6wwa891/wCFnhueG4v7gurKpYtnivQwEcJKVsS2u1iZc3Q85uvFOiv8S4tTXaLMPkkeleu2/wAS/Ct1KIob8Mx6DbXzvFp+ir4kMcjk2Cvhj7V7J4f8DeBNRZLjTJd0gHTdzX1OcYXBRhB1efSNlZfmYU5S1sdZeJofjO1ezEquwGcgc1x3hbw/rXhXX76BAz2RB2ntisLxLpGteCteGpaYZDabscelereHdc/tvw+tyykS7CHz64ryaqnhMNejLnpT79GaL3nrueTaOxfx7cO33i9e5RuFtVZugXJrwzRz/wAV9Pn++f517lGgktFQ9CuKjPfjp+iClszn9Q8f+HNMm8m7vQj+mKqf8LR8Jf8AQRX8qi1b4WeH9ZuTPdLIXPoazz8FPC39yX86ypwyjkXPKd/RDftOhpSfFDwiY2H9oK3HTHWvGtG1/SYviq2qXEgFiXyGr0q++DfhW2s5ZiXTYpOWPFeOaTo2k3Pjs6ZcPiy37d2a+gyill7pVnQcmuXW/by8zKo53Vz6A/4Wl4R/6CI/75pf+Fo+Ev8AoIr+VZC/BjwpKodBIynkENxTv+FKeFv7kv514fJkv80/uRpep5HQ2Pj7w5qAb7Neh9vXiqlz8Q/CTF7We9Rg3ysrLxUOmfC3QNKDi3WT5+uTXP8Aiz4UeH4NDv7+MyJOiFwSe9TRpZVOty807O1tgbnY6nw/pHh25d7/AESRSC3Ozsa6K+hkkspVV+dprxr4QfbE8J6qLBi1wCdgHrVS5ufimY5BIknl810VsqnPFTh7Ze40veeolP3b2I/DcMh+LGwN82/rXvuxsffNfJunSeJh4p3Wob+0s/rXqOg3HxJa9Yaisnk7e/rXfneWyqSjP2kVaK3er9Cac7dDm/Dcbv8AGO7VeSJD/OvocfdA9q+evh75p+LF2bj/AF285zX0L3rzeI9K9OPaKLo7MNvuaTb/ALRp1FfPXNRu05+8a8g+K6O8wVSSQa9g715B8XVmiRpohzXsZE/9tiZ1fhLuneMbHw/4BRYJAb4Lwg65rm/CfhPUfGesHV9c3iAksoauY8BfYzqIuddLG3DcA9K+jdIu7G8sE/s4qIFGAFr18xn/AGY5xor3p7y7eSM4e/a5LYWEOm28dtAuI1HFecfGXIsLdk/1g+7XqPQgV5f8X/8AVWX1/rXjZPJyx8JM0qfCzgdJ8S/ECDSvKsraRrYDg46CpvBFzf3fi1JNUUrLu5zXsXg8N/wiCcR/cOPy715Tqk72WrzXSgb1Y/dr6KjjIYmVajGlGL2utzJxtZ3Ol+J2r213cwaTC26dsAAV2HgbT5dM8PRwT5Vs5xXlGk+JfD7a7HqWsn9/GcgNXoB+MXhNePtL8ei15+NweJjh4YWjTbS1bt1LjJX5mzm/ifC2nypftuKBsk10Xhf4j+GX0e0tzdhJ9u0oR3qvf/EPwL4jtzZ38oaM93HSm6T8N/BeoMl9pjFlU5+U0p+y+qRpY6E4uOzS0DXmvFnoEt7ENLkvFYeWIywP4V4hpEK+O/GUkjtuFu+R+FdN8TfET6d4fOl6QSSq7X29hXK/BS5jtZL+7k5fnNaZfhJ4fAVcXH4npHuKcryUT17xJrNt4e8OStLIFkSHCD1IryrwBpMvizXn128BIjfK7hWB4w8ZnWfE4huw/wBijfDemK9J0T4i+CtH0uOC0lEQCjcAOSatYLEYHB2pwcqlTqui7C5lKWr0R6HdW8F7btBKAUPUYrw/xDa3vgHxX/atqGFkWyQBxXc/8Lj8J5/4+X/75qOfx54L8VQnT7qUMj9N4rhwFLGYOT9pRk4PdW6FzcZbPU6Pwz4qsPEenRTQSjzGHzL71leI/hppniSYy3U8qtnPy1yll4R1HSdXW40IsbBmyMdMV6Je6/b6Dp0L6nJiVh096xrQeGrqeAn8Wy6ryGnzL3jzzUfg14a0zT5ri4vpECqSCxxzXmOhWHhZtZlg1O8YwK2FYntW98RvEWr6hdRlmkXTmPHParGhJ8NX09P7UYi4x8x96+qw0sVSwvPiJym5fy6tGEuVysjSPhr4VuARqxX23Uf8Iv8ACv8A6DB/76o+yfCLP+uf86Psnwiz/rn/ADrm9pU/nrfch2XkPj8I/C+Zgsescn1euj0z4PeHFkjvrK8lkjPKkHINeZeMLTwILJT4fncT59a9t+GkbJ4KswzMx9Sa5czq4nD4ZVqdaertaSsVBJys0jhvjbAtno2mWycoHAyfrV7xXry6P8NdKt45CGnjVeDVb48cWOnf9dB/OrniPwyuvfDTTbiMEvbxCT8qnDyp/VcI623Mwd+aVjY8D6Quj+Dp70gb54i/4YrzPQ/BEHjLxHdmecxDeTwa7z4c+J11vw1d6dMcNaxlMGvLo38RReIro6CGOHP3a6MFDERrYlc3LPu9rdBSasj0k/AjRcf8f1zn/PvWN4r+Dml6N4dub+3vJmkhGcP0Nc3L478fQXy6bI0guWOFU5p3iW7+Ii6JIdXWX7H/AB5zW9KjmkKsPaYlWb2vuvITcGnZHonwc40HH93ivTj65xivMfgyPM8OeYDj1FL8T/FupaRALfSgxkf5Ttr57G4SeKzOdKG7ZrGSjC5jfFfxq7smjaU5aZztfbW98LvBKaDYDUZxm5nXJz1rnfhv4Curu6Ot66rM8nzKGr2YbIIOMLFGv6CtsyxVPDUFgMM7/wAz7sUIuT52PzySK5zxh4ah17SJhsBuAnyHHermn+J9M1K6e3tpt8iHBFbBIHNeHF1cNUUrWaNdJI8P8FeK7nwrqx0HVCwQvhS1e2xSJNGskbAowyCK4fx34Eg1y3e+tE2368qR3rkfC3irWvDV2NL1wPgNgZr3MTQp5lT+sYfSf2o9/NGSbg7PY3/GuueMLG6ZdFtmkjzwQtcd/wAJd8Tv+fKT/vmvc7a9jurQXMYOwjNcnqPxS8N6Zcvb3FwwkQ4IFRgsU3H2UMLGbW+mo5R68x5v/wAJf8Tv+fKT/vij/hL/AInf8+Un/fFd3/wuXwqP+W8n5UD4y+FD/wAt5Pyrv9piP+gBfcyLL+Y4X/hLvicR/wAeUn/fNdL4V13xte3G3VrZ44j1JGK6rSviLoOszCK0nLOeADXUSPmFjjgiuDF47lXs54WMG/LUuMeqkeFIAPi5Dg5JfmveT1FeCxf8lej/AN+vej1FTnu9H/CgpdRaKKK8A1CiiigBD1xXnfxC8QpbSpoaHMl0uOD613d/I0WnzSr95UJFeT2Hh+88S+JP7Uuw37hvlz7V62VUqfM61V6R/PoZzb2RwWkeFo7PxxFpuojME5zz716ja/D650PxVHdadN5enjqpNcreyeZ8WLRX6IwFeieN9H8QatsTSLnykxzzivezDGVpVKUZTUVOOt9v+HM4RVmb+o3WlzQG3vZI3j75NYn/AAk+g6JavDbD5Bn7tcXa/DvxOR/pd4W55+aumsvASR2Ey3nzOVNeRLD4OiuWVXmXZF3k+h5RaeNdP0/xfcX0ke5C5wK9S0r4q6bqk0dvDCwZsDpXmWg+GdOl8az215GGiD9DXt9l4N0CzKSW9kisOjV6mczy+LipRblbR9CKfObqEPGr+oBpSoPHNHyqu3O0AVy3i3xxYeFrbLMsspHCqc4r5WjRqV5qFJXbN20ldmH8WPESaf4eksbZ83T8YB5rzhPBk1r4DTXzG32wnd71u6J4fv8Ax54hGu3QZbMNnY1eyvpdpJp32AxL9n27dtfSSxscrpww9N3le8v8jHl522znPh5qqal4Vtlkf/SFGGBPNdaqjGOa8SvYdS+HfiZr9Vd9OduEHQCvVdA8S2PiCwS4hkVWI5QtyK8zMsG4v6xS1hLW/byLhLo9za2iue8bKP8AhD9Q/wCuZroQfTn3rA8aLv8ACN+ucExmuDCP9/D1X5ly2Z578Bj/AMSzURj/AJadfxr1m/kaOylZQCdpr5/+HHjW28NaBqdky7rpmJTHrzXo3grVb/VdBuZr4MC2SN1e9neBqvF1MRJWjdfMypyXKkcH4bmcfFnfgZ39K99ycfdr5/8ADf8AyVf/ALaV9BHoax4ht7Wn/hQ6WzPnvw6ssPxhu5RxmQ/zr6CHIBz2rxGezbTPHc14/wAis/U1reKPFd9bXFottIQjkAkGunMsNLHVaXs39lfgTB8qdz1qiqmnStLpdvI5+ZowTVuvl5R5W0bhXm3xfh8vw010B0ODXpDdDXL/ABA0o6x4SubVVyx5FduWVVSxdOUtromavFnKeAvDGk6/4FilmiDSPkZHY1mXmj+JfCU/nWbubFT0HpU3wo8SQWDHw3cMI5UY7dxxzXrdxDDdW7RygNGRzXrY3F1sHjJwqLmhJ3s+z7GcYqUVY5zwl4qi16ARNxOg+bNcp8ZebG3VP9Zj5awLnWItJ8Zi20s8GTDBaufGTUkFvpIRgZHIyK2w2B9lmFKcFZS1SFKV4NM5fSdI+Ik2lF7GaQWuOme1O8H2l4/iVbXXCWBbDZr2rw/5tt4NtmC/N5O6vLrRmk8bAv1Mn9a7KWYzxPtocsVa+qWpLhaxN8Svh/plgE1eKMrbL98LWn4N8C+DPEmgx3kcG9s4YA4xXo+taTbazoz2F0wWN1714He6xP8ADrWZLGxfzLUknCmubA4jE4/C/V6dRqpHbXdFSSjK7Who/Erwl4a0+2htdFULfhsFVOa9K+HGhDSPB0CMD5zoc5rgPBGiN4o8R/23fTqIj8wjY17hFGkMapGuEA4xXPnGKnToRwTm5Natvv2HTjd8xwi+Dtw1O6vhuDoxXNcF8HLaKXxJq1qwzEGYAV7jqBJ0q7z/AM8n/ka8T+DH/I3at/vNVYLE1K2AxLk9kgkkpRPTZ/h54cuHdpbMEv15rxPxl4S0ex+JFlpNiCttKV8xQenrX0Fr2ow6To9zeSOqmNMgE9TXi3hO0l8c+LzrRjIWF85NVkuJxEI1MRUm+SKa36vYVRLRJHoafCPwkqKDZEkDBO7rXmHizwZpOk+JUi08FQG6A19DYA/hrxTxdIsvjcQJ9/f0qMmx2KqVpc9RtWfUdSMUtj1XwwjxeHrZCMlVwK5bV/CWoa/rSSXzE2ytkCux0SNotJhVuoFX8nI6c140cVOhWnOnu76mnLdanjPxisLax07T7WBAqHg/nWh4W+F3hbUNEgubiPzZXXLYbpVb45DdZ2aLw56fnXJ6H4J8c3mnxzWN+0cJHA3Yr6jDc88spv2/s3d6vqYP43pc9P8A+FQ+D/8An0b/AL7o/wCFQ+D/APn0b/vuuD/4V98Rf+go3/fym/8ACAfEUH/kJt/38rD2dT/oPX3sq6/lLfxG+Hvh7QNEF3p6iOYNjaWya9C+Gjl/BVoW68ivEvFfhTxbpVis+sXjSQ56Fs17f8N2VvBdmV6VObJrLYc1X2nvbhT+PaxxHx5/48NP/wB8fzrvPD6xr8PLZZJF2G0OSTx0NcL8eCP7Lsx/EWGPzrGtPBXjPU/C9mLW+ZIJIx8u7Hy0QoQrZXQVSooJSe4NtTdkV/CsyWV7qwsgSGLZIrX+G8xXxHKSAQznOa6Lwb4E/wCEX8O3rajiS4dCS3XtXM+AvJHiO6Z51Rd5xk101sRSxEMR7PVJJX7kpNWuaHj7TLmz8b2WsxJ+4jILHHFdlqHizQrrQXkvgsluy/OnXms7xT4ksdS3+H4l8y4mXajjnBryWx0288OeLU0nXNz2UrfxdMZrmw+FWMoQ9v7sqa0XVryKcuV6dTsfBeoT3fiVl0aNo9MLdMcV6rd6Dp99KJLiEO49aj0TStIsbdW0uOMRkcFDmtbg8kV4uPxvta3NSTjbTz+ZpGNlqMSNYo1iRcIowAK4H4n+Lh4c0wW0TZlnUqAOvNdfrWsW+iadJczuBhSVBPU14ho1hd/ErxVJc3gYWtu2V3dMV05RhIzk8VX/AIcNX5sVSX2VudN8JPD1zD5up3e7998wzXe6t4t07RtRjsrpiskmMc1qWdvbWFpHbRMirGoHUVxHxG8KQ63YPqEFyq3MC5X5vSodenj8dzYjSL0Xl2CzjHQ72KWOaNZYiGVhkEVnaj4d03VLhZ7qANIO9eYfDPx/suP7B1J8zBtqsTXsfy8Y5z3rnxmFrZfXcG7dn3Q4yU0VxbRWmnyQxDCBDgfhXzz4S8Oad4j+Juo22pAyRKzMF9TX0XP/AMe0v+6f5V826P4ng8GfEbUb24jMyszLhevNetkPtp0sQqN+dx07kVbXVz14/CTwkW/48j/31UN18JvCSWcrCzKkKTuz0rDPx30oHH9nz1HP8ddLktpEXT5tzKQM0Rwue3Xxff8A8EOakcj4d0i00zxkI7UttD8CvoZ/+PQf7o/lXzr4V1aXWfFouI7N0RnznFfRT5+ygbT93+lLiJTVamqm9go2s7HhcX/JXo/9/wDrXvR6ivBYsf8AC3o+f4696PUVjnu9H/Ch0uotFFFfPmoUUUUAIyq6FWGVPUVAlvDbxP5EapkdqsUYpptaAeCXcV5/wtqBhbN5XmfexXvLEgjGKgNjZmcTG3j83+9t5qzXoY/HLFKmuW3KrERjy3E+b2pkhbyZPXBqSivPTLPCNIivP+E9nL27BPMPzY969wjIS1V8dFzimrZWiymRYIw5/i281YwMYxxXoY/HLFSi+W1lYiMeU8u8ZeLtRmn/ALMsLd1ZuN4FZug/C7Uby9S/1y5M0bHdsY5r1ptPsml8xraMv/eK81ZAxjHSuiObOjR9lho8vd9Rezu7sr2NjbadbLb2kSxxr2FWcUUV48pOTuzQp6jplnqluYLyFZEIwMjpXkGrfDzW9D1CXUtMu2FqpLCJW7V7XTWUOpVwCp6g13YLMa2EbUNU909iZQUjzTwZ4/vL+5Gm3dm4dDt3kV0Pj2edfD0sUCFmkXGBXQRabYQyeZFaxK/94LzViWGKVcSorD0IqquLoPERrU6dkuglF2s2eFfDf4bST3r6jqKfuw33WFezXFlb2mmSR2sKoAvRavxRRxLtiRUX0Ap5AIwRkU8fmdXGVvaT27BGCirI+f8Aw5FeD4rbntmEW/72K9/b7y1XSxs0m81LeMSf3gvNWqMxx6xk4yUbWVghHlPO/ih4fmvtKFzYLiaPk7e9edeF5L/xVew6dc2ro1ueXI64r6HkRJEKuoZT1Bqtb6dY20hkt7aKNz1ZVwa6sLnDo4Z0XG7Wz7ClTu7klrB9nsooP7ihanoorxG23dmgU10VlIcAqexp1FIDxX4gfDe9j1STxBo0hjdTuKp1rEt/iZruj6c1ncWcskpG3dg19BuiupVwCp6g1QfRtKkOXsYGPulfQ0M6i6UaWLp86jt3MnT1vF2PFPhj4bvNZ1u51nUYmUHLKGFTeI/D1/4h8WRq6N5ED8enFe4W9tb2ybbeJI19FGKBbwBywiQMe+KU8+qPESrqNtLJdg9krWPOvE/xAPhHT7exTT3mbyxGCAeOK4rw/eXmq3j68bRoxGxO0ivc7nS9Pu2BubSGQjoWXNOi06xhiMcVtEiHqoXioo5nh6NFxhS957u+4ODb3PHtV8Ua54sYWthFJblflyBitPQPhZLIPtGtSedIw/i5r06HT7GBt0VtEjeoWrXNRUziUIezwseRfj941T1vI8e1D4dazpN6b7TLsrbJz5SntW14f8e3JnXTrqzcuvy7yK9GI3AhsEHtVVdNsFk8xbWIP67eaiWZqvT5cTDma2ewcln7o27kMuiXLheWhcgfga8W+ES3Vrr2s3c1uyqpY8jGete7bV2bcDbjGKrw2NnBv8m3jTd97auM1lhceqGHq0eW/Pb5WHKN2n2PFNV/tzx/r0ltEJIbNHwy9ARXoMFlbfD7wpJLb2++RVywA6murgs7W3YtBCiMepUYqSaCKeIxzIroeoYZFa4jM/aqNJRtTXTuJQtr1PGo/js6qRJo0hb6Gsvw1Ff+NPGw1eWBoYy2cEV7T/wj+i/9A62/791btrGztBi2t44v91cV1PNcJShL6rR5ZSVr3uLkk/iZLEgijWMdhinnqKWivn7mp5B8bY7l00/7PC0mDzgdOa7TwE9yfDsCzLtwo4I6V0lxaW1zj7RCkmOm4ZqSKKOJAsaBV9AK9Krj1PBQwvL8L3IUbS5hs0jRRM+AdozivMdf+LzaJqLWv9lPLj+IA16kQCMEZFUJtG0qd981jA7erJmscFWw1OTeIp8y9bDkm9meBeIfGuq/EK4i02GxkggJ/umvbfBelPo3hm2s5PvKMmtCHR9Lt3Dw2UCMO6pitAAAYHSuvMMyp16UaFCHJBa28yYwad2zx347RXEljYeRC0nz87RnFeieEDN/wh+m+Yu1xAMg9q17m0troAXMKSAdNwzUsaJHGERQqAYAHasK2OVTB08Ny/C27+o1G0nIp6krSaPcqRljGeBXznY+CvEOpeIJxau8KmQnPTivpkgEEEZFRR28ETFo40Vj3ArXLs1ngYTUI3chTgpbnHeFfAi6O6XN+wmul6OecVb8aeDovE1pmPCXSj5XrqxnvS/SuZ4+u66r83vIrlVrHivhfVNZ8I6z/ZN9HJNHuwGPNevte7dOa72Hhd22nSWNnLL5slvG0n94rzU5jQpsKjb6VeNxkMVNVOSz6+YoxcdD5q+IfjHWPEGqtaw20yQQtjCqfmrsPC3iV4NDjtbSwaGfbh229a9XOj6WXLGyg3HqdlSR6Zp8RzHaxL9Fr1Kuc4aeHjQjRso+ZCpu97nmkejeINaZmS4eIH3xQnw78RlyJNVJjPVd9eqJGkYxGoX6Cnc+tcTzistKaSXoV7NdTw/xB8NLjQoTrFo2+5Q5+WvRPh/qd3qWgKb1GWZTzurqpI0lQpKoZT1BFNgghgXbDGqL6KMVOJzSeKoezrK8k9H5dhqCi7odIm+F0/vAivMdE+E0Nr4sutW1MpcwyElY25wTXqNJg556Vy4fG1sPGcaTtzKzG4p7mD/whnhvOf7Ktvyo/wCEO8Ng/wDIKts/7tb+B6UmB6VH1uv/ADv72PlXY8/8T+JNG8APEtro6vI/OY06Vhz/ABob7OTHpEgZhxwa9Su9Nsbwg3VrFKR03rmoP7D0j/nwt/8AvivQo4vBKC9tScpdXchxlfRnjng/Sr3xB4xXX5YmiTdnBFe6Eciobe1trZdtvCkY9FGKnrnzDHPF1FK1klZLyKhHlQUUUV55QUUUUAf/2Q==")
	resp,_ := http.PostForm("https://api-cn.faceplusplus.com/imagepp/v1/recognizetext",value)
	bs,_:=ioutil.ReadAll(resp.Body)
	log.Println(string(bs))*/

	//baidu
	value := make(url.Values)
	value.Add("language_type", "CHN_ENG")
	value.Add("Content-Type", "application/x-www-form-urlencoded")
	value.Add("image", "iVBORw0KGgoAAAANSUhEUgAAAKAAAAARBAMAAACod7rOAAAAG1BMVEX///8AAP9fX/+fn/9/f//f3/8/P/+/v/8fH/8EVqDgAAACKUlEQVQ4jdWTS3PaMBRGv9qS8dIqGLI0kCYsDU2cLIE+l3gy43rp4Dy6NJ1J66UDofCze69kaOAXtPIMMnocnXuvDPwHrQ985H6DEtYYmOpfRDMkIXAKJAVgZzJDMuMXmpNKRUopGv9JD3CruFV/gbyyBpo5FUBk8gUaKALNmeOB/jsBA703dJhLwDniouY09kDBhNmhoW24za5qPlPPwORmjYkZfudFjjYUPh18BOyqjnnZAXdN+tCGGyC2lBoSoQXkesrLHW24CFnS0qec7wxXpntcZeUFAa1Q53RSR26Aew+9t5K+r4FySUPp+NBQNjFUJ+SzWZULAooKsvXKUGqge1/HE/JO6fV0yAtdoeXsACja9xFSUvA/l+6ITnuBVdVRqLffbjs6h8NtyIuTa2PoxclV5RYjCoysxzEHswOm+VQDL/2ktKcEXITDYm94tkSHDZePPdrkfXBrQycK15KsxZqXxkTdA3t9jAgIu6qLYrcoYmEMWxQyc0Ve53AHTKNwdE2L49ApjoCf+pT4tBFYQQ2Uz2ZzOuDcnNq8TtzsgCbkdtsJXZaba8cDIPqyQ8CMLqgBfs1/c8hyLjYMdHQR7SNDOFf3dD2ogvhOwNc5RD8NEIuAPgMGNroriO0XNAYhyrzA3WJ2DBQMdNtFOcbw1/ThLGBDnbMaSF+jc/KUGWB5znOD90uu6o/VZH0BA+Ritxk4VCEB6brKvKHuLiHNjdqae/Dvtz8b53kRh1rTrAAAAABJRU5ErkJggg==")
	resp, _ := http.PostForm("https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic?access_token=24.cb9afce813b0dc1449529649bfc2a3dc.2592000.1570863615.282335-17235293", value)
	bs, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(bs))
}

// 企查查接口
func Qichacha() {

}

func get() {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("https://www.qichacha.com/search?key=小米")
	req.Header.Add("Cookie", "acw_tc=6f7e78a715668012034494032efb97f66ec87cb68fb44b21e12d5f8b0d; zg_did=%7B%22did%22%3A%20%2216ccca116221a2-0f050a5165c1-38607701-13c680-16ccca116235f5%22%7D; UM_distinctid=16ccca11b11b0-0284524a81e822-38607701-13c680-16ccca11b12b6e; QCCSESSID=4ddso719atrpgkbmukneu4spe7; hasShow=1; CNZZDATA1254842228=1950090108-1566797068-https%253A%252F%252Fwww.baidu.com%252F%7C1568626442; Hm_lvt_3456bee468c83cc63fb5147f119f1075=1568202279,1568254024,1568614161,1568628815; Hm_lpvt_3456bee468c83cc63fb5147f119f1075=1568628815; zg_de1d1a35bfa24ce29bbf2c7eb17e6c4f=%7B%22sid%22%3A%201568625694403%2C%22updated%22%3A%201568628815114%2C%22info%22%3A%201568614160398%2C%22superProperty%22%3A%20%22%7B%7D%22%2C%22platform%22%3A%20%22%7B%7D%22%2C%22utm%22%3A%20%22%7B%7D%22%2C%22referrerDomain%22%3A%20%22www.qichacha.com%22%2C%22zs%22%3A%200%2C%22sc%22%3A%200%2C%22cuid%22%3A%20%221d0083f0fc0802d50c303cfd3dedfa79%22%7D")
	req.Header.Add("Host", "www.qichacha.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Cache-Control", "max-age=0")

	var client fasthttp.Client

	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(resp.Body()))
}

func download() {
	imgPath := "."
	imgUrl := "http://data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAKAAAAARBAMAAACod7rOAAAAG1BMVEX///8AAP/f3/+/v/9/f/9fX/8fH/8/P/+fn/83sGquAAACUElEQVQ4jdWTT1PbMBDFH7Js+Wi50HCUBya5OtOOuSadOuTohEI4Ovw/GmggRwHF6cfurhwSQr9Aq8x4HHn3p/d2V8C/v8TXkxTK0psX06NjplrrVNGjRX+V3vG0LoAfwCuHhxXnaL3gKOCAfsCt5lU2wMHts+ytgf45amyl/OUQKlU290pJwHPcJ4ZPKBgYbXUsAoqaIUmXyrwlEBnEbA3E6TvgoyFgWBEwLMNSVug7KXocLZRTKGPI6m8g1IEDRs3GCqgqtiwLsqlsYEkOMHdliubKKexaFum7U46WwJHWxxQvm7O3pxEDjQMSgBRmRkyrrpkA7XGTW4o4dkBxSWHt3geFA6MqZzlkyyKmjL4D9isHPOX36f4ZuDMjy5kiOnaWu65Dl2YTeHWBIwf0t52blcJ6wu3PqYW7cnRAzd9Bp24URknnpQxS6vs95fUS2izegOFDx0A9EyKpV0BXw8ybGaqh3xOxpA7HOMYwWCpUC3sm7gDJwpEQdQUMMuz9Gn2maZiwTgcM0iZMFqRQxH5BmukDzeIbsL2wrzUVL7Eq/QDMM3l9mlFceCbOGTiuoV542gj7PUhzHA4MAWk8GNhYbrWo7Sxu5jRuAL/R2LibcmJAXNG6qyGHXG1ZYpaYHO1dEJAGeL5WSGfe3bgQPBDwfQ2xBNZUEXGdejtkWVwUrl2/Pw2Q71/WY2ztXQg2wEDJwKCVZj08fckfJwUrJIebwPDebd6EBQHRd0X0p8WwuztLsbB01RK+0gR80paAPpmZe/rnM7g3fPqVxX+x/gC4nIDbXy+SzAAAAABJRU5ErkJggg=="

	fileName := "tmp.png"

	res, err := http.Get(imgUrl)
	if err != nil {
		fmt.Println("A error occurred:", err)
		return
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)

	file, err := os.Create(imgPath + fileName)
	if err != nil {
		//panic(err)
		fmt.Errorf("err:%s", err)
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	written, _ := io.Copy(writer, reader)
	fmt.Printf("Total length: %d", written)
}

func initHttpProxyServer() {
	ips := []string{"39.107.243.195"}
	sshd.InitServer(ips)
	pids := sshd.ExecuteBatch(`ps -ef|grep httpproxy | grep -v grep | awk -F " " '{print $2}' `)

	for ip, _ := range serverMap {
		sshd.Execute(ip, "kill -9 "+pids[ip])
	}

	sshd.UploadBatch("./pkg/httpproxy/httpproxy", "/var/tmp/")
	sshd.ExecuteBatch(`echo "" >> /var/tmp/startHttpProxy.sh`)
	sshd.ExecuteBatch(`echo "/var/tmp/httpproxy 1 >> /dev/null 2>&1 &" > startHttpProxy.sh`)
	sshd.ExecuteBatch(`chmod +x /var/tmp/httpproxy`)
	sshd.ExecuteBatch(`chmod +x startHttpProxy.sh`)
	sshd.ExecuteBatch(`./startHttpProxy.sh`)
	sshd.ExecuteBatch(`ps -ef|grep httpproxy | grep -v grep | awk -F " " '{print $2}' `)
}

func httpd(ip, urlStr string) string {
	proxyUrl, err := url.Parse("http://" + ip + ":8888")
	myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}

	res, err := myClient.Get(urlStr)

	if err != nil {
		log.Println(err)
	}

	bs, _ := ioutil.ReadAll(res.Body)
	//log.Println(string(bs))
	return string(bs)
}

func t() {
	s := `<div class="market-header" data-s-efc25ec2="">
                <div class="match-module-container" id="match-module-container-quote-new-chart-pankou-0.725358720181755" data-s-0bd93592="" style="position: relative;">
        <!--308-->
        <div class="c-gap-bottom-middle c-gap-bottom-lh" data-s-d84eef82="">
            <div class="pankou-container" data-s-d84eef82="" style="max-height: 66px;">
                <div class="pankou-fold-box" data-s-d84eef82="">
                    <div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">今开
                            <!--311-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-green); font-weight: 500;">35.80</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">最高
                            <!--312-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-red); font-weight: 500;">36.40</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">成交量
                            <!--313-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-flat); font-weight: 500;">5.00万手</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">
                            <span class="capitalization-popover">
        
        <!--885-->
            <!--886-->
                                <span class="cos-color-text-tiny" data-s-d84eef82="" style="cursor: pointer;">总市值
                                
                                
                                <svg style="transform: rotate(-180deg); margin-bottom: 2px;" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="8px" height="6px" viewBox="0 0 12 8" version="1.1" data-s-d84eef82="">
                                    <g id="控件" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" data-s-d84eef82="">
                                        <g id="盘口数据/默认展示3行" transform="translate(-181.000000, -130.000000)" fill="#858585" data-s-d84eef82="">
                                            <g id="编组" transform="translate(24.000000, 88.000000)" data-s-d84eef82="">
                                                <g id="拖动/横滑组件备份-2" transform="translate(118.000000, 26.000000)" data-s-d84eef82="">
                                                    <g id="编组-2" transform="translate(39.000000, 16.000000)" data-s-d84eef82="">
                                                        
                                                        <path d="M6.8479983,1.35679729 L10.0437507,6.47000106 C10.336461,6.93833759 10.1940878,7.55528797 9.7257513,7.8479983 C9.5668197,7.94733056 9.38317208,8 9.19575236,8 L2.80424764,8 C2.25196289,8 1.80424764,7.55228475 1.80424764,7 C1.80424764,6.81258028 1.85691709,6.62893266 1.95624934,6.47000106 L5.1520017,1.35679729 C5.44471203,0.888460755 6.06166241,0.74608759 6.52999894,1.03879792 C6.65876807,1.11927863 6.7675176,1.22802815 6.8479983,1.35679729 Z" id="三角形备份" transform="translate(6.000000, 4.000000) scale(1, -1) rotate(-180.000000) translate(-6.000000, -4.000000) " data-s-d84eef82=""></path>
                                                    </g>
                                                </g>
                                            </g>
                                        </g>
                                    </g>
                                </svg>
                            <!--886-->
            
        <!--885-->
    
    <!--314-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-flat); font-weight: 500;">567.48亿</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">昨收
                            <!--315-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-flat); font-weight: 500;">35.84</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">最低
                            <!--316-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-green); font-weight: 500;">35.67</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">成交额
                            <!--317-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-flat); font-weight: 500;">1.81亿</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">总股本
                            <!--318-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-flat); font-weight: 500;">15.70亿</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">换手率
                            <!--319-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-flat); font-weight: 500;">0.32%</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">振幅
                            <!--324-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-flat); font-weight: 500;">2.04%</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">量比
                            <!--325-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-green); font-weight: 500;">0.77</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">市盈(TTM)
                            <!--326-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-flat); font-weight: 500;">20.28</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">涨跌幅
                            <!--327-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-red); font-weight: 500;">+0.84%</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">均价
                            <!--328-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-red); font-weight: 500;">36.15</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">委比
                            <!--329-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-flat); font-weight: 500;">-50.77%</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">市盈(静)
                            <!--330-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-flat); font-weight: 500;">22.00</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">流通值
                            <!--331-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-flat); font-weight: 500;">567.47亿</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">涨停
                            <!--332-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-red); font-weight: 500;">39.42</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">内盘
                            <!--333-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-green); font-weight: 500;">2.69万</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">市净率
                            <!--334-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-flat); font-weight: 500;">3.92</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">流通股
                            <!--335-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-flat); font-weight: 500;">15.70亿</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">跌停
                            <!--336-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-green); font-weight: 500;">32.26</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">外盘
                            <!--337-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-red); font-weight: 500;">2.31万</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">市销率
                            <!--338-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-flat); font-weight: 500;">--</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">52周高
                            <!--339-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-flat); font-weight: 500;">37.89</div>
                    </div><div class="pankou-item" data-s-d84eef82="" style="width: 25%;">
                        <div class="key" data-s-d84eef82="">
                            <span class="cos-color-text-tiny" data-s-d84eef82="">52周低
                            <!--340-->
                        </div>
                        <div class="value harmony-os-medium" data-s-d84eef82="" style="color: var(--stock-color-flat); font-weight: 500;">21.14</div>
                    </div><!--310-->
                </div>
            </div><!--309-->
            <div class="flod-box-btn cos-color-text" data-s-d84eef82="">
                <i class="cos-icon cos-icon-down"></i>
            </div><!--341-->
            <!--343-->
        </div>
    <!--308-->
    <object tabindex="-1" type="text/html" aria-hidden="true" data="about:blank" style="display: block; position: absolute; top: 0px; left: 0px; width: 100%; height: 100%; border: none; padding: 0px; margin: 0px; opacity: 0; z-index: -1000; pointer-events: none;"></object></div>
            </div>
            <div class="chart-container" data-s-efc25ec2="">
                <div class="fac-chart-container" data-s-9dda213a="">
        <div class="fac-chart-container-header" data-s-9dda213a="">
            <div class="fac-chart-container-tabs" data-s-9dda213a="">
                <div class="match-module-container" id="match-module-container-quote-new-chart-tabs-0.09361391250578255" data-s-0bd93592="" style="position: relative;">
        <!--348-->
        <div class="
            capsule-tabs
            flex-left
            middle-screen-box
            
        " data-s-2529b0ff="">
            <div class="left-tabs" data-s-2529b0ff="">
                <!--349-->
            </div>
            <div class="public-tab-wraper" data-s-2529b0ff="">
                <div class="
                        public-tab
                        c-border-raidus-middle
                        active-tag
                        
                    " data-s-2529b0ff="">
                    分时
                    <!--351-->
                </div><div class="
                        public-tab
                        c-border-raidus-middle
                        
                        
                    " data-s-2529b0ff="">
                    五日
                    <!--352-->
                </div><div data-s-2529b0ff="" class="
                        public-tab
                        c-border-raidus-middle
                        
                        
                    ">
                    日K
                    <!--353-->
                </div><div data-s-2529b0ff="" class="
                        public-tab
                        c-border-raidus-middle
                        
                        
                    ">
                    周K
                    <!--354-->
                </div><div data-s-2529b0ff="" class="
                        public-tab
                        c-border-raidus-middle
                        
                        
                    ">
                    月K
                    <!--355-->
                </div><div data-s-2529b0ff="" class="
                        public-tab
                        c-border-raidus-middle
                        
                        
                    ">
                    季K
                    <!--356-->
                </div><div data-s-2529b0ff="" class="
                        public-tab
                        c-border-raidus-middle
                        
                        
                    ">
                    年K
                    <!--357-->
                </div><div data-s-2529b0ff="" class="
                        public-tab
                        c-border-raidus-middle
                        
                         public-tab-no-padding
                    ">
                    
                    <div class="chart-dropdown-container cos-color-text" data-s-2e87cd65="">
        <div class="entry" data-s-2e87cd65="">
            <div class="entry-btn " data-s-2e87cd65="">
                <span data-s-2e87cd65="" style="color: rgb(0, 3, 17);">更多
                <i class="cos-icon cos-icon-down"></i>
            </div>
            <div class="dropdown-menu-list-wrapper" data-s-2e87cd65="" style="display: none;">
                <div class="dropdown-menu-list" data-s-2e87cd65="">
                    <div class="dropdown-menu-item" data-s-2e87cd65="">
                        1分
                    </div><div class="dropdown-menu-item" data-s-2e87cd65="">
                        5分
                    </div><div class="dropdown-menu-item" data-s-2e87cd65="">
                        15分
                    </div><div class="dropdown-menu-item" data-s-2e87cd65="">
                        30分
                    </div><div class="dropdown-menu-item" data-s-2e87cd65="">
                        60分
                    </div><!--361-->
                </div>
            </div>
        </div>
    </div><!--358-->
                </div><!--350-->
            </div>
            <div class="right-tabs" data-s-2529b0ff="">
                
                <!--362-->
                
                <div class="right-button" data-s-2529b0ff="">
                    <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg" data-s-2529b0ff="">
                        <path fill-rule="evenodd" clip-rule="evenodd" d="M2 4.66667C2 3.19391 3.19391 2 4.66667 2H7.41667V3.33333H4.66667C3.93029 3.33333 3.33333 3.93029 3.33333 4.66667V7.41667H2V4.66667ZM10.5833 3.33333V2H13.3333C14.8061 2 16 3.19391 16 4.66667V7.41667H14.6667V4.66667C14.6667 3.93029 14.0697 3.33333 13.3333 3.33333H10.5833ZM3.33333 10.5833V13.3333C3.33333 14.0697 3.93029 14.6667 4.66667 14.6667H7.41667V16H4.66667C3.19391 16 2 14.8061 2 13.3333V10.5833H3.33333ZM14.6667 10.5833H16V13.3333C16 14.8061 14.8061 16 13.3333 16H10.5833V14.6667H13.3333C14.0697 14.6667 14.6667 14.0697 14.6667 13.3333V10.5833Z" fill="black" data-s-2529b0ff=""></path>
                    </svg><!--364-->

                    <span class="button-text entry" data-s-2529b0ff="">超级视图<!--365-->
                </div><!--363-->
            </div>
        </div>
    <!--348-->
    <object tabindex="-1" type="text/html" aria-hidden="true" data="about:blank" style="display: block; position: absolute; top: 0px; left: 0px; width: 100%; height: 100%; border: none; padding: 0px; margin: 0px; opacity: 0; z-index: -1000; pointer-events: none;"></object></div>
            </div>

            
            <!--366-->
        </div>

        <div class="flex-container" data-s-9dda213a="">
            <div class="left-chart-container" data-s-9dda213a="" style="width: calc(100% - 237px);">
                <div class="main-chart-container" data-s-9dda213a="">
                    <div data-s-0bd93592="" class="match-module-container" id="match-module-container-fac-market-chart-new-chart-0.06155975233600486" style="position: relative;">
        <!--368-->
                        <div class="chart-container" data-s-9dda213a="">
                            <div class="chart" data-s-9dda213a="">
                                
                                <div class="svg" data-s-9dda213a=""><!--887--><svg version="1.1" xmlns="http://www.w3.org/2000/svg" width="524px" height="324px"><text x="0" y="244.24" fill="#3D404D" text-anchor="start" style="font-size: 12px; font-weight: 400; font-family: ;">09:30</text>,<text x="264.1833333333333" y="244.24" fill="#3D404D" text-anchor="middle" style="font-size: 12px; font-weight: 400; font-family: ;">11:30/13:00</text>,<text x="524" y="244.24" fill="#3D404D" text-anchor="end" style="font-size: 12px; font-weight: 400; font-family: ;">15:00</text>
        <mask id="higherMask">  
            <rect width="100%" height="100%" fill="white"></rect>  
            <rect y="186.63428571428707" width="100%" height="100%" fill="black"></rect>  
        </mask>
        <mask id="lowerMask">  
            <rect width="100%" height="100%" fill="black"></rect>  
            <rect y="186.63428571428707" width="100%" height="100%" fill="white"></rect>  
        </mask>
        
            <polyline points="0,199.792380952382 2.183333333333333,166.89714285714354 4.366666666666666,186.63428571428474 6.549999999999999,216.24000000000004 8.733333333333333,186.63428571428474 10.916666666666666,173.47619047618983 13.099999999999998,166.89714285714354 15.283333333333331,163.6076190476204 17.466666666666665,160.3180952380949 19.65,166.89714285714354 21.833333333333332,173.47619047618983 24.016666666666666,157.02857142857178 26.199999999999996,160.3180952380949 28.38333333333333,163.6076190476204 30.566666666666663,163.6076190476204 32.75,153.73904761904865 34.93333333333333,173.47619047618983 37.11666666666666,166.89714285714354 39.3,153.73904761904865 41.48333333333333,160.3180952380949 43.666666666666664,166.89714285714354 45.849999999999994,193.21333333333337 48.03333333333333,186.63428571428474 50.21666666666666,186.63428571428474 52.39999999999999,183.34476190476158 54.58333333333333,134.0019047619051 56.76666666666666,140.58095238095373 58.949999999999996,140.58095238095373 61.133333333333326,143.87047619047686 63.31666666666666,134.0019047619051 65.5,117.55428571428706 67.68333333333332,120.84380952381021 69.86666666666666,101.10666666666668 72.05,120.84380952381021 74.23333333333332,114.26476190476158 76.41666666666666,107.68571428571529 78.6,87.94857142857177 80.78333333333333,97.81714285714354 82.96666666666665,94.52761904762039 85.14999999999999,94.52761904762039 87.33333333333333,74.79047619047687 89.51666666666665,71.50095238095372 91.69999999999999,64.9219047619051 93.88333333333333,45.18476190476157 96.06666666666666,41.89523809523843 98.24999999999999,38.605714285715294 100.43333333333332,58.34285714285882 102.61666666666666,55.05333333333334 104.79999999999998,58.34285714285882 106.98333333333332,58.34285714285882 109.16666666666666,61.63238095238196 111.35,58.34285714285882 113.53333333333332,51.763809523810195 115.71666666666665,51.763809523810195 117.89999999999999,45.18476190476157 120.08333333333333,58.34285714285882 122.26666666666665,45.18476190476157 124.44999999999999,35.31619047619215 126.63333333333333,38.605714285715294 128.81666666666666,41.89523809523843 131,28.737142857143528 133.1833333333333,22.158095238094905 135.36666666666665,25.447619047620385 137.54999999999998,12.28952380952548 139.73333333333332,32.02666666666667 141.91666666666666,22.158095238094905 144.1,35.31619047619215 146.28333333333333,41.89523809523843 148.46666666666664,32.02666666666667 150.64999999999998,15.579047619048621 152.83333333333331,25.447619047620385 155.01666666666665,41.89523809523843 157.2,32.02666666666667 159.38333333333333,64.9219047619051 161.56666666666666,51.763809523810195 163.74999999999997,71.50095238095372 165.9333333333333,55.05333333333334 168.11666666666665,55.05333333333334 170.29999999999998,68.21142857142824 172.48333333333332,71.50095238095372 174.66666666666666,68.21142857142824 176.85,71.50095238095372 179.0333333333333,61.63238095238196 181.21666666666664,78.08000000000001 183.39999999999998,64.9219047619051 185.58333333333331,74.79047619047687 187.76666666666665,84.65904761904864 189.95,97.81714285714354 192.13333333333333,94.52761904762039 194.31666666666663,87.94857142857177 196.49999999999997,68.21142857142824 198.6833333333333,68.21142857142824 200.86666666666665,71.50095238095372 203.04999999999998,58.34285714285882 205.23333333333332,51.763809523810195 207.41666666666666,48.474285714287056 209.59999999999997,41.89523809523843 211.7833333333333,45.18476190476157 213.96666666666664,25.447619047620385 216.14999999999998,35.31619047619215 218.33333333333331,38.605714285715294 220.51666666666665,45.18476190476157 222.7,48.474285714287056 224.88333333333333,51.763809523810195 227.06666666666663,38.605714285715294 229.24999999999997,35.31619047619215 231.4333333333333,22.158095238094905 233.61666666666665,32.02666666666667 235.79999999999998,28.737142857143528 237.98333333333332,32.02666666666667 240.16666666666666,22.158095238094905 242.34999999999997,22.158095238094905 244.5333333333333,12.28952380952548 246.71666666666664,9 248.89999999999998,32.02666666666667 251.08333333333331,28.737142857143528 253.26666666666665,28.737142857143528 255.45,28.737142857143528 257.6333333333333,12.28952380952548 259.81666666666666,25.447619047620385 262,28.737142857143528 264.1833333333333,45.18476190476157 266.3666666666666,38.605714285715294 268.54999999999995,22.158095238094905 270.7333333333333,28.737142857143528 272.91666666666663,18.868571428571762 275.09999999999997,22.158095238094905 277.2833333333333,25.447619047620385 279.46666666666664,35.31619047619215 281.65,35.31619047619215 283.8333333333333,32.02666666666667 286.01666666666665,45.18476190476157 288.2,58.34285714285882 290.3833333333333,61.63238095238196 292.56666666666666,68.21142857142824 294.75,58.34285714285882 296.9333333333333,48.474285714287056 299.1166666666666,55.05333333333334 301.29999999999995,61.63238095238196 303.4833333333333,64.9219047619051 305.66666666666663,55.05333333333334 307.84999999999997,51.763809523810195 310.0333333333333,55.05333333333334 312.21666666666664,58.34285714285882 314.4,55.05333333333334 316.5833333333333,58.34285714285882 318.76666666666665,64.9219047619051 320.95,61.63238095238196 323.1333333333333,61.63238095238196 325.31666666666666,58.34285714285882 327.49999999999994,64.9219047619051 329.6833333333333,64.9219047619051 331.8666666666666,64.9219047619051 334.04999999999995,68.21142857142824 336.2333333333333,71.50095238095372 338.41666666666663,74.79047619047687 340.59999999999997,71.50095238095372 342.7833333333333,81.36952380952549 344.96666666666664,81.36952380952549 347.15,84.65904761904864 349.3333333333333,81.36952380952549 351.51666666666665,78.08000000000001 353.7,81.36952380952549 355.8833333333333,81.36952380952549 358.0666666666666,84.65904761904864 360.24999999999994,81.36952380952549 362.4333333333333,78.08000000000001 364.6166666666666,81.36952380952549 366.79999999999995,78.08000000000001 368.9833333333333,71.50095238095372 371.16666666666663,87.94857142857177 373.34999999999997,87.94857142857177 375.5333333333333,94.52761904762039 377.71666666666664,84.65904761904864 379.9,87.94857142857177 382.0833333333333,97.81714285714354 384.26666666666665,91.23809523809491 386.45,87.94857142857177 388.63333333333327,87.94857142857177 390.8166666666666,84.65904761904864 392.99999999999994,74.79047619047687 395.1833333333333,91.23809523809491 397.3666666666666,81.36952380952549 399.54999999999995,71.50095238095372 401.7333333333333,94.52761904762039 403.91666666666663,91.23809523809491 406.09999999999997,91.23809523809491 408.2833333333333,91.23809523809491 410.46666666666664,87.94857142857177 412.65,87.94857142857177 414.8333333333333,81.36952380952549 417.01666666666665,81.36952380952549 419.19999999999993,87.94857142857177 421.38333333333327,78.08000000000001 423.5666666666666,71.50095238095372 425.74999999999994,61.63238095238196 427.9333333333333,58.34285714285882 430.1166666666666,74.79047619047687 432.29999999999995,74.79047619047687 434.4833333333333,71.50095238095372 436.66666666666663,71.50095238095372 438.84999999999997,55.05333333333334 441.0333333333333,64.9219047619051 443.21666666666664,71.50095238095372 445.4,74.79047619047687 447.5833333333333,74.79047619047687 449.76666666666665,74.79047619047687 451.94999999999993,71.50095238095372 454.13333333333327,64.9219047619051 456.3166666666666,61.63238095238196 458.49999999999994,71.50095238095372 460.6833333333333,71.50095238095372 462.8666666666666,81.36952380952549 465.04999999999995,87.94857142857177 467.2333333333333,74.79047619047687 469.41666666666663,81.36952380952549 471.59999999999997,84.65904761904864 473.7833333333333,81.36952380952549 475.96666666666664,81.36952380952549 478.15,81.36952380952549 480.3333333333333,87.94857142857177 482.5166666666666,87.94857142857177 484.69999999999993,87.94857142857177 486.88333333333327,97.81714285714354 489.0666666666666,91.23809523809491 491.24999999999994,87.94857142857177 493.4333333333333,97.81714285714354 495.6166666666666,94.52761904762039 497.79999999999995,74.79047619047687 499.9833333333333,91.23809523809491 502.16666666666663,91.23809523809491 504.34999999999997,97.81714285714354 506.5333333333333,97.81714285714354 508.71666666666664,91.23809523809491 510.9,91.23809523809491 513.0833333333333,91.23809523809491 515.2666666666667,81.36952380952549 517.4499999999999,101.10666666666668 519.6333333333333,87.94857142857177 521.8166666666666,87.94857142857177 524,87.94857142857177" style="fill: none;stroke: #037B66;" stroke-width="1.2" mask="url(#lowerMask)"></polyline>
            <polyline points="0,199.792380952382 2.183333333333333,166.89714285714354 4.366666666666666,186.63428571428474 6.549999999999999,216.24000000000004 8.733333333333333,186.63428571428474 10.916666666666666,173.47619047618983 13.099999999999998,166.89714285714354 15.283333333333331,163.6076190476204 17.466666666666665,160.3180952380949 19.65,166.89714285714354 21.833333333333332,173.47619047618983 24.016666666666666,157.02857142857178 26.199999999999996,160.3180952380949 28.38333333333333,163.6076190476204 30.566666666666663,163.6076190476204 32.75,153.73904761904865 34.93333333333333,173.47619047618983 37.11666666666666,166.89714285714354 39.3,153.73904761904865 41.48333333333333,160.3180952380949 43.666666666666664,166.89714285714354 45.849999999999994,193.21333333333337 48.03333333333333,186.63428571428474 50.21666666666666,186.63428571428474 52.39999999999999,183.34476190476158 54.58333333333333,134.0019047619051 56.76666666666666,140.58095238095373 58.949999999999996,140.58095238095373 61.133333333333326,143.87047619047686 63.31666666666666,134.0019047619051 65.5,117.55428571428706 67.68333333333332,120.84380952381021 69.86666666666666,101.10666666666668 72.05,120.84380952381021 74.23333333333332,114.26476190476158 76.41666666666666,107.68571428571529 78.6,87.94857142857177 80.78333333333333,97.81714285714354 82.96666666666665,94.52761904762039 85.14999999999999,94.52761904762039 87.33333333333333,74.79047619047687 89.51666666666665,71.50095238095372 91.69999999999999,64.9219047619051 93.88333333333333,45.18476190476157 96.06666666666666,41.89523809523843 98.24999999999999,38.605714285715294 100.43333333333332,58.34285714285882 102.61666666666666,55.05333333333334 104.79999999999998,58.34285714285882 106.98333333333332,58.34285714285882 109.16666666666666,61.63238095238196 111.35,58.34285714285882 113.53333333333332,51.763809523810195 115.71666666666665,51.763809523810195 117.89999999999999,45.18476190476157 120.08333333333333,58.34285714285882 122.26666666666665,45.18476190476157 124.44999999999999,35.31619047619215 126.63333333333333,38.605714285715294 128.81666666666666,41.89523809523843 131,28.737142857143528 133.1833333333333,22.158095238094905 135.36666666666665,25.447619047620385 137.54999999999998,12.28952380952548 139.73333333333332,32.02666666666667 141.91666666666666,22.158095238094905 144.1,35.31619047619215 146.28333333333333,41.89523809523843 148.46666666666664,32.02666666666667 150.64999999999998,15.579047619048621 152.83333333333331,25.447619047620385 155.01666666666665,41.89523809523843 157.2,32.02666666666667 159.38333333333333,64.9219047619051 161.56666666666666,51.763809523810195 163.74999999999997,71.50095238095372 165.9333333333333,55.05333333333334 168.11666666666665,55.05333333333334 170.29999999999998,68.21142857142824 172.48333333333332,71.50095238095372 174.66666666666666,68.21142857142824 176.85,71.50095238095372 179.0333333333333,61.63238095238196 181.21666666666664,78.08000000000001 183.39999999999998,64.9219047619051 185.58333333333331,74.79047619047687 187.76666666666665,84.65904761904864 189.95,97.81714285714354 192.13333333333333,94.52761904762039 194.31666666666663,87.94857142857177 196.49999999999997,68.21142857142824 198.6833333333333,68.21142857142824 200.86666666666665,71.50095238095372 203.04999999999998,58.34285714285882 205.23333333333332,51.763809523810195 207.41666666666666,48.474285714287056 209.59999999999997,41.89523809523843 211.7833333333333,45.18476190476157 213.96666666666664,25.447619047620385 216.14999999999998,35.31619047619215 218.33333333333331,38.605714285715294 220.51666666666665,45.18476190476157 222.7,48.474285714287056 224.88333333333333,51.763809523810195 227.06666666666663,38.605714285715294 229.24999999999997,35.31619047619215 231.4333333333333,22.158095238094905 233.61666666666665,32.02666666666667 235.79999999999998,28.737142857143528 237.98333333333332,32.02666666666667 240.16666666666666,22.158095238094905 242.34999999999997,22.158095238094905 244.5333333333333,12.28952380952548 246.71666666666664,9 248.89999999999998,32.02666666666667 251.08333333333331,28.737142857143528 253.26666666666665,28.737142857143528 255.45,28.737142857143528 257.6333333333333,12.28952380952548 259.81666666666666,25.447619047620385 262,28.737142857143528 264.1833333333333,45.18476190476157 266.3666666666666,38.605714285715294 268.54999999999995,22.158095238094905 270.7333333333333,28.737142857143528 272.91666666666663,18.868571428571762 275.09999999999997,22.158095238094905 277.2833333333333,25.447619047620385 279.46666666666664,35.31619047619215 281.65,35.31619047619215 283.8333333333333,32.02666666666667 286.01666666666665,45.18476190476157 288.2,58.34285714285882 290.3833333333333,61.63238095238196 292.56666666666666,68.21142857142824 294.75,58.34285714285882 296.9333333333333,48.474285714287056 299.1166666666666,55.05333333333334 301.29999999999995,61.63238095238196 303.4833333333333,64.9219047619051 305.66666666666663,55.05333333333334 307.84999999999997,51.763809523810195 310.0333333333333,55.05333333333334 312.21666666666664,58.34285714285882 314.4,55.05333333333334 316.5833333333333,58.34285714285882 318.76666666666665,64.9219047619051 320.95,61.63238095238196 323.1333333333333,61.63238095238196 325.31666666666666,58.34285714285882 327.49999999999994,64.9219047619051 329.6833333333333,64.9219047619051 331.8666666666666,64.9219047619051 334.04999999999995,68.21142857142824 336.2333333333333,71.50095238095372 338.41666666666663,74.79047619047687 340.59999999999997,71.50095238095372 342.7833333333333,81.36952380952549 344.96666666666664,81.36952380952549 347.15,84.65904761904864 349.3333333333333,81.36952380952549 351.51666666666665,78.08000000000001 353.7,81.36952380952549 355.8833333333333,81.36952380952549 358.0666666666666,84.65904761904864 360.24999999999994,81.36952380952549 362.4333333333333,78.08000000000001 364.6166666666666,81.36952380952549 366.79999999999995,78.08000000000001 368.9833333333333,71.50095238095372 371.16666666666663,87.94857142857177 373.34999999999997,87.94857142857177 375.5333333333333,94.52761904762039 377.71666666666664,84.65904761904864 379.9,87.94857142857177 382.0833333333333,97.81714285714354 384.26666666666665,91.23809523809491 386.45,87.94857142857177 388.63333333333327,87.94857142857177 390.8166666666666,84.65904761904864 392.99999999999994,74.79047619047687 395.1833333333333,91.23809523809491 397.3666666666666,81.36952380952549 399.54999999999995,71.50095238095372 401.7333333333333,94.52761904762039 403.91666666666663,91.23809523809491 406.09999999999997,91.23809523809491 408.2833333333333,91.23809523809491 410.46666666666664,87.94857142857177 412.65,87.94857142857177 414.8333333333333,81.36952380952549 417.01666666666665,81.36952380952549 419.19999999999993,87.94857142857177 421.38333333333327,78.08000000000001 423.5666666666666,71.50095238095372 425.74999999999994,61.63238095238196 427.9333333333333,58.34285714285882 430.1166666666666,74.79047619047687 432.29999999999995,74.79047619047687 434.4833333333333,71.50095238095372 436.66666666666663,71.50095238095372 438.84999999999997,55.05333333333334 441.0333333333333,64.9219047619051 443.21666666666664,71.50095238095372 445.4,74.79047619047687 447.5833333333333,74.79047619047687 449.76666666666665,74.79047619047687 451.94999999999993,71.50095238095372 454.13333333333327,64.9219047619051 456.3166666666666,61.63238095238196 458.49999999999994,71.50095238095372 460.6833333333333,71.50095238095372 462.8666666666666,81.36952380952549 465.04999999999995,87.94857142857177 467.2333333333333,74.79047619047687 469.41666666666663,81.36952380952549 471.59999999999997,84.65904761904864 473.7833333333333,81.36952380952549 475.96666666666664,81.36952380952549 478.15,81.36952380952549 480.3333333333333,87.94857142857177 482.5166666666666,87.94857142857177 484.69999999999993,87.94857142857177 486.88333333333327,97.81714285714354 489.0666666666666,91.23809523809491 491.24999999999994,87.94857142857177 493.4333333333333,97.81714285714354 495.6166666666666,94.52761904762039 497.79999999999995,74.79047619047687 499.9833333333333,91.23809523809491 502.16666666666663,91.23809523809491 504.34999999999997,97.81714285714354 506.5333333333333,97.81714285714354 508.71666666666664,91.23809523809491 510.9,91.23809523809491 513.0833333333333,91.23809523809491 515.2666666666667,81.36952380952549 517.4499999999999,101.10666666666668 519.6333333333333,87.94857142857177 521.8166666666666,87.94857142857177 524,87.94857142857177" style="fill: none;stroke: #D60A22;" stroke-width="1.2" mask="url(#higherMask)"></polyline>
            <defs>
                <linearGradient id="gradient-up" x1="0%" y1="100%" x2="0%" y2="85.71428571428636%">
                    <stop offset="0%" stop-color="#037B66" stop-opacity="0.8"></stop>
                    <stop offset="100%" stop-color="#037B66" stop-opacity="0"></stop>
                </linearGradient>
            </defs>
                <polygon points="0,9 0,199.792380952382 2.183333333333333,166.89714285714354 4.366666666666666,186.63428571428474 6.549999999999999,216.24000000000004 8.733333333333333,186.63428571428474 10.916666666666666,173.47619047618983 13.099999999999998,166.89714285714354 15.283333333333331,163.6076190476204 17.466666666666665,160.3180952380949 19.65,166.89714285714354 21.833333333333332,173.47619047618983 24.016666666666666,157.02857142857178 26.199999999999996,160.3180952380949 28.38333333333333,163.6076190476204 30.566666666666663,163.6076190476204 32.75,153.73904761904865 34.93333333333333,173.47619047618983 37.11666666666666,166.89714285714354 39.3,153.73904761904865 41.48333333333333,160.3180952380949 43.666666666666664,166.89714285714354 45.849999999999994,193.21333333333337 48.03333333333333,186.63428571428474 50.21666666666666,186.63428571428474 52.39999999999999,183.34476190476158 54.58333333333333,134.0019047619051 56.76666666666666,140.58095238095373 58.949999999999996,140.58095238095373 61.133333333333326,143.87047619047686 63.31666666666666,134.0019047619051 65.5,117.55428571428706 67.68333333333332,120.84380952381021 69.86666666666666,101.10666666666668 72.05,120.84380952381021 74.23333333333332,114.26476190476158 76.41666666666666,107.68571428571529 78.6,87.94857142857177 80.78333333333333,97.81714285714354 82.96666666666665,94.52761904762039 85.14999999999999,94.52761904762039 87.33333333333333,74.79047619047687 89.51666666666665,71.50095238095372 91.69999999999999,64.9219047619051 93.88333333333333,45.18476190476157 96.06666666666666,41.89523809523843 98.24999999999999,38.605714285715294 100.43333333333332,58.34285714285882 102.61666666666666,55.05333333333334 104.79999999999998,58.34285714285882 106.98333333333332,58.34285714285882 109.16666666666666,61.63238095238196 111.35,58.34285714285882 113.53333333333332,51.763809523810195 115.71666666666665,51.763809523810195 117.89999999999999,45.18476190476157 120.08333333333333,58.34285714285882 122.26666666666665,45.18476190476157 124.44999999999999,35.31619047619215 126.63333333333333,38.605714285715294 128.81666666666666,41.89523809523843 131,28.737142857143528 133.1833333333333,22.158095238094905 135.36666666666665,25.447619047620385 137.54999999999998,12.28952380952548 139.73333333333332,32.02666666666667 141.91666666666666,22.158095238094905 144.1,35.31619047619215 146.28333333333333,41.89523809523843 148.46666666666664,32.02666666666667 150.64999999999998,15.579047619048621 152.83333333333331,25.447619047620385 155.01666666666665,41.89523809523843 157.2,32.02666666666667 159.38333333333333,64.9219047619051 161.56666666666666,51.763809523810195 163.74999999999997,71.50095238095372 165.9333333333333,55.05333333333334 168.11666666666665,55.05333333333334 170.29999999999998,68.21142857142824 172.48333333333332,71.50095238095372 174.66666666666666,68.21142857142824 176.85,71.50095238095372 179.0333333333333,61.63238095238196 181.21666666666664,78.08000000000001 183.39999999999998,64.9219047619051 185.58333333333331,74.79047619047687 187.76666666666665,84.65904761904864 189.95,97.81714285714354 192.13333333333333,94.52761904762039 194.31666666666663,87.94857142857177 196.49999999999997,68.21142857142824 198.6833333333333,68.21142857142824 200.86666666666665,71.50095238095372 203.04999999999998,58.34285714285882 205.23333333333332,51.763809523810195 207.41666666666666,48.474285714287056 209.59999999999997,41.89523809523843 211.7833333333333,45.18476190476157 213.96666666666664,25.447619047620385 216.14999999999998,35.31619047619215 218.33333333333331,38.605714285715294 220.51666666666665,45.18476190476157 222.7,48.474285714287056 224.88333333333333,51.763809523810195 227.06666666666663,38.605714285715294 229.24999999999997,35.31619047619215 231.4333333333333,22.158095238094905 233.61666666666665,32.02666666666667 235.79999999999998,28.737142857143528 237.98333333333332,32.02666666666667 240.16666666666666,22.158095238094905 242.34999999999997,22.158095238094905 244.5333333333333,12.28952380952548 246.71666666666664,9 248.89999999999998,32.02666666666667 251.08333333333331,28.737142857143528 253.26666666666665,28.737142857143528 255.45,28.737142857143528 257.6333333333333,12.28952380952548 259.81666666666666,25.447619047620385 262,28.737142857143528 264.1833333333333,45.18476190476157 266.3666666666666,38.605714285715294 268.54999999999995,22.158095238094905 270.7333333333333,28.737142857143528 272.91666666666663,18.868571428571762 275.09999999999997,22.158095238094905 277.2833333333333,25.447619047620385 279.46666666666664,35.31619047619215 281.65,35.31619047619215 283.8333333333333,32.02666666666667 286.01666666666665,45.18476190476157 288.2,58.34285714285882 290.3833333333333,61.63238095238196 292.56666666666666,68.21142857142824 294.75,58.34285714285882 296.9333333333333,48.474285714287056 299.1166666666666,55.05333333333334 301.29999999999995,61.63238095238196 303.4833333333333,64.9219047619051 305.66666666666663,55.05333333333334 307.84999999999997,51.763809523810195 310.0333333333333,55.05333333333334 312.21666666666664,58.34285714285882 314.4,55.05333333333334 316.5833333333333,58.34285714285882 318.76666666666665,64.9219047619051 320.95,61.63238095238196 323.1333333333333,61.63238095238196 325.31666666666666,58.34285714285882 327.49999999999994,64.9219047619051 329.6833333333333,64.9219047619051 331.8666666666666,64.9219047619051 334.04999999999995,68.21142857142824 336.2333333333333,71.50095238095372 338.41666666666663,74.79047619047687 340.59999999999997,71.50095238095372 342.7833333333333,81.36952380952549 344.96666666666664,81.36952380952549 347.15,84.65904761904864 349.3333333333333,81.36952380952549 351.51666666666665,78.08000000000001 353.7,81.36952380952549 355.8833333333333,81.36952380952549 358.0666666666666,84.65904761904864 360.24999999999994,81.36952380952549 362.4333333333333,78.08000000000001 364.6166666666666,81.36952380952549 366.79999999999995,78.08000000000001 368.9833333333333,71.50095238095372 371.16666666666663,87.94857142857177 373.34999999999997,87.94857142857177 375.5333333333333,94.52761904762039 377.71666666666664,84.65904761904864 379.9,87.94857142857177 382.0833333333333,97.81714285714354 384.26666666666665,91.23809523809491 386.45,87.94857142857177 388.63333333333327,87.94857142857177 390.8166666666666,84.65904761904864 392.99999999999994,74.79047619047687 395.1833333333333,91.23809523809491 397.3666666666666,81.36952380952549 399.54999999999995,71.50095238095372 401.7333333333333,94.52761904762039 403.91666666666663,91.23809523809491 406.09999999999997,91.23809523809491 408.2833333333333,91.23809523809491 410.46666666666664,87.94857142857177 412.65,87.94857142857177 414.8333333333333,81.36952380952549 417.01666666666665,81.36952380952549 419.19999999999993,87.94857142857177 421.38333333333327,78.08000000000001 423.5666666666666,71.50095238095372 425.74999999999994,61.63238095238196 427.9333333333333,58.34285714285882 430.1166666666666,74.79047619047687 432.29999999999995,74.79047619047687 434.4833333333333,71.50095238095372 436.66666666666663,71.50095238095372 438.84999999999997,55.05333333333334 441.0333333333333,64.9219047619051 443.21666666666664,71.50095238095372 445.4,74.79047619047687 447.5833333333333,74.79047619047687 449.76666666666665,74.79047619047687 451.94999999999993,71.50095238095372 454.13333333333327,64.9219047619051 456.3166666666666,61.63238095238196 458.49999999999994,71.50095238095372 460.6833333333333,71.50095238095372 462.8666666666666,81.36952380952549 465.04999999999995,87.94857142857177 467.2333333333333,74.79047619047687 469.41666666666663,81.36952380952549 471.59999999999997,84.65904761904864 473.7833333333333,81.36952380952549 475.96666666666664,81.36952380952549 478.15,81.36952380952549 480.3333333333333,87.94857142857177 482.5166666666666,87.94857142857177 484.69999999999993,87.94857142857177 486.88333333333327,97.81714285714354 489.0666666666666,91.23809523809491 491.24999999999994,87.94857142857177 493.4333333333333,97.81714285714354 495.6166666666666,94.52761904762039 497.79999999999995,74.79047619047687 499.9833333333333,91.23809523809491 502.16666666666663,91.23809523809491 504.34999999999997,97.81714285714354 506.5333333333333,97.81714285714354 508.71666666666664,91.23809523809491 510.9,91.23809523809491 513.0833333333333,91.23809523809491 515.2666666666667,81.36952380952549 517.4499999999999,101.10666666666668 519.6333333333333,87.94857142857177 521.8166666666666,87.94857142857177 524,87.94857142857177 524,9" fill="url(#gradient-up)" fill-opacity="0.2" mask="url(#lowerMask)"></polygon>
            
            <defs>
                <linearGradient id="gradient-down" x1="0%" y1="0%" x2="0%" y2="77.82785038305602%">
                    <stop offset="0%" stop-color="#D60A22" stop-opacity="0.8"></stop>
                    <stop offset="100%" stop-color="#D60A22" stop-opacity="0"></stop>
                </linearGradient>
            </defs>
                <polygon points="0,237.24 0,199.792380952382 2.183333333333333,166.89714285714354 4.366666666666666,186.63428571428474 6.549999999999999,216.24000000000004 8.733333333333333,186.63428571428474 10.916666666666666,173.47619047618983 13.099999999999998,166.89714285714354 15.283333333333331,163.6076190476204 17.466666666666665,160.3180952380949 19.65,166.89714285714354 21.833333333333332,173.47619047618983 24.016666666666666,157.02857142857178 26.199999999999996,160.3180952380949 28.38333333333333,163.6076190476204 30.566666666666663,163.6076190476204 32.75,153.73904761904865 34.93333333333333,173.47619047618983 37.11666666666666,166.89714285714354 39.3,153.73904761904865 41.48333333333333,160.3180952380949 43.666666666666664,166.89714285714354 45.849999999999994,193.21333333333337 48.03333333333333,186.63428571428474 50.21666666666666,186.63428571428474 52.39999999999999,183.34476190476158 54.58333333333333,134.0019047619051 56.76666666666666,140.58095238095373 58.949999999999996,140.58095238095373 61.133333333333326,143.87047619047686 63.31666666666666,134.0019047619051 65.5,117.55428571428706 67.68333333333332,120.84380952381021 69.86666666666666,101.10666666666668 72.05,120.84380952381021 74.23333333333332,114.26476190476158 76.41666666666666,107.68571428571529 78.6,87.94857142857177 80.78333333333333,97.81714285714354 82.96666666666665,94.52761904762039 85.14999999999999,94.52761904762039 87.33333333333333,74.79047619047687 89.51666666666665,71.50095238095372 91.69999999999999,64.9219047619051 93.88333333333333,45.18476190476157 96.06666666666666,41.89523809523843 98.24999999999999,38.605714285715294 100.43333333333332,58.34285714285882 102.61666666666666,55.05333333333334 104.79999999999998,58.34285714285882 106.98333333333332,58.34285714285882 109.16666666666666,61.63238095238196 111.35,58.34285714285882 113.53333333333332,51.763809523810195 115.71666666666665,51.763809523810195 117.89999999999999,45.18476190476157 120.08333333333333,58.34285714285882 122.26666666666665,45.18476190476157 124.44999999999999,35.31619047619215 126.63333333333333,38.605714285715294 128.81666666666666,41.89523809523843 131,28.737142857143528 133.1833333333333,22.158095238094905 135.36666666666665,25.447619047620385 137.54999999999998,12.28952380952548 139.73333333333332,32.02666666666667 141.91666666666666,22.158095238094905 144.1,35.31619047619215 146.28333333333333,41.89523809523843 148.46666666666664,32.02666666666667 150.64999999999998,15.579047619048621 152.83333333333331,25.447619047620385 155.01666666666665,41.89523809523843 157.2,32.02666666666667 159.38333333333333,64.9219047619051 161.56666666666666,51.763809523810195 163.74999999999997,71.50095238095372 165.9333333333333,55.05333333333334 168.11666666666665,55.05333333333334 170.29999999999998,68.21142857142824 172.48333333333332,71.50095238095372 174.66666666666666,68.21142857142824 176.85,71.50095238095372 179.0333333333333,61.63238095238196 181.21666666666664,78.08000000000001 183.39999999999998,64.9219047619051 185.58333333333331,74.79047619047687 187.76666666666665,84.65904761904864 189.95,97.81714285714354 192.13333333333333,94.52761904762039 194.31666666666663,87.94857142857177 196.49999999999997,68.21142857142824 198.6833333333333,68.21142857142824 200.86666666666665,71.50095238095372 203.04999999999998,58.34285714285882 205.23333333333332,51.763809523810195 207.41666666666666,48.474285714287056 209.59999999999997,41.89523809523843 211.7833333333333,45.18476190476157 213.96666666666664,25.447619047620385 216.14999999999998,35.31619047619215 218.33333333333331,38.605714285715294 220.51666666666665,45.18476190476157 222.7,48.474285714287056 224.88333333333333,51.763809523810195 227.06666666666663,38.605714285715294 229.24999999999997,35.31619047619215 231.4333333333333,22.158095238094905 233.61666666666665,32.02666666666667 235.79999999999998,28.737142857143528 237.98333333333332,32.02666666666667 240.16666666666666,22.158095238094905 242.34999999999997,22.158095238094905 244.5333333333333,12.28952380952548 246.71666666666664,9 248.89999999999998,32.02666666666667 251.08333333333331,28.737142857143528 253.26666666666665,28.737142857143528 255.45,28.737142857143528 257.6333333333333,12.28952380952548 259.81666666666666,25.447619047620385 262,28.737142857143528 264.1833333333333,45.18476190476157 266.3666666666666,38.605714285715294 268.54999999999995,22.158095238094905 270.7333333333333,28.737142857143528 272.91666666666663,18.868571428571762 275.09999999999997,22.158095238094905 277.2833333333333,25.447619047620385 279.46666666666664,35.31619047619215 281.65,35.31619047619215 283.8333333333333,32.02666666666667 286.01666666666665,45.18476190476157 288.2,58.34285714285882 290.3833333333333,61.63238095238196 292.56666666666666,68.21142857142824 294.75,58.34285714285882 296.9333333333333,48.474285714287056 299.1166666666666,55.05333333333334 301.29999999999995,61.63238095238196 303.4833333333333,64.9219047619051 305.66666666666663,55.05333333333334 307.84999999999997,51.763809523810195 310.0333333333333,55.05333333333334 312.21666666666664,58.34285714285882 314.4,55.05333333333334 316.5833333333333,58.34285714285882 318.76666666666665,64.9219047619051 320.95,61.63238095238196 323.1333333333333,61.63238095238196 325.31666666666666,58.34285714285882 327.49999999999994,64.9219047619051 329.6833333333333,64.9219047619051 331.8666666666666,64.9219047619051 334.04999999999995,68.21142857142824 336.2333333333333,71.50095238095372 338.41666666666663,74.79047619047687 340.59999999999997,71.50095238095372 342.7833333333333,81.36952380952549 344.96666666666664,81.36952380952549 347.15,84.65904761904864 349.3333333333333,81.36952380952549 351.51666666666665,78.08000000000001 353.7,81.36952380952549 355.8833333333333,81.36952380952549 358.0666666666666,84.65904761904864 360.24999999999994,81.36952380952549 362.4333333333333,78.08000000000001 364.6166666666666,81.36952380952549 366.79999999999995,78.08000000000001 368.9833333333333,71.50095238095372 371.16666666666663,87.94857142857177 373.34999999999997,87.94857142857177 375.5333333333333,94.52761904762039 377.71666666666664,84.65904761904864 379.9,87.94857142857177 382.0833333333333,97.81714285714354 384.26666666666665,91.23809523809491 386.45,87.94857142857177 388.63333333333327,87.94857142857177 390.8166666666666,84.65904761904864 392.99999999999994,74.79047619047687 395.1833333333333,91.23809523809491 397.3666666666666,81.36952380952549 399.54999999999995,71.50095238095372 401.7333333333333,94.52761904762039 403.91666666666663,91.23809523809491 406.09999999999997,91.23809523809491 408.2833333333333,91.23809523809491 410.46666666666664,87.94857142857177 412.65,87.94857142857177 414.8333333333333,81.36952380952549 417.01666666666665,81.36952380952549 419.19999999999993,87.94857142857177 421.38333333333327,78.08000000000001 423.5666666666666,71.50095238095372 425.74999999999994,61.63238095238196 427.9333333333333,58.34285714285882 430.1166666666666,74.79047619047687 432.29999999999995,74.79047619047687 434.4833333333333,71.50095238095372 436.66666666666663,71.50095238095372 438.84999999999997,55.05333333333334 441.0333333333333,64.9219047619051 443.21666666666664,71.50095238095372 445.4,74.79047619047687 447.5833333333333,74.79047619047687 449.76666666666665,74.79047619047687 451.94999999999993,71.50095238095372 454.13333333333327,64.9219047619051 456.3166666666666,61.63238095238196 458.49999999999994,71.50095238095372 460.6833333333333,71.50095238095372 462.8666666666666,81.36952380952549 465.04999999999995,87.94857142857177 467.2333333333333,74.79047619047687 469.41666666666663,81.36952380952549 471.59999999999997,84.65904761904864 473.7833333333333,81.36952380952549 475.96666666666664,81.36952380952549 478.15,81.36952380952549 480.3333333333333,87.94857142857177 482.5166666666666,87.94857142857177 484.69999999999993,87.94857142857177 486.88333333333327,97.81714285714354 489.0666666666666,91.23809523809491 491.24999999999994,87.94857142857177 493.4333333333333,97.81714285714354 495.6166666666666,94.52761904762039 497.79999999999995,74.79047619047687 499.9833333333333,91.23809523809491 502.16666666666663,91.23809523809491 504.34999999999997,97.81714285714354 506.5333333333333,97.81714285714354 508.71666666666664,91.23809523809491 510.9,91.23809523809491 513.0833333333333,91.23809523809491 515.2666666666667,81.36952380952549 517.4499999999999,101.10666666666668 519.6333333333333,87.94857142857177 521.8166666666666,87.94857142857177 524,87.94857142857177 524,237.24" fill="url(#gradient-down)" fill-opacity="0.2" mask="url(#higherMask)"></polygon>
            
        
        <polyline points="0,200.792380952382 2.183333333333333,189.9238095238102 4.366666666666666,189.9238095238102 6.549999999999999,189.9238095238102 8.733333333333333,189.9238095238102 10.916666666666666,189.9238095238102 13.099999999999998,189.9238095238102 15.283333333333331,189.9238095238102 17.466666666666665,186.63428571428474 19.65,186.63428571428474 21.833333333333332,186.63428571428474 24.016666666666666,183.34476190476158 26.199999999999996,183.34476190476158 28.38333333333333,183.34476190476158 30.566666666666663,180.05523809523845 32.75,180.05523809523845 34.93333333333333,180.05523809523845 37.11666666666666,180.05523809523845 39.3,180.05523809523845 41.48333333333333,176.76571428571532 43.666666666666664,176.76571428571532 45.849999999999994,176.76571428571532 48.03333333333333,176.76571428571532 50.21666666666666,176.76571428571532 52.39999999999999,176.76571428571532 54.58333333333333,173.47619047618983 56.76666666666666,173.47619047618983 58.949999999999996,173.47619047618983 61.133333333333326,173.47619047618983 63.31666666666666,173.47619047618983 65.5,170.1866666666667 67.68333333333332,170.1866666666667 69.86666666666666,166.89714285714354 72.05,163.6076190476204 74.23333333333332,163.6076190476204 76.41666666666666,163.6076190476204 78.6,157.02857142857178 80.78333333333333,157.02857142857178 82.96666666666665,153.73904761904865 85.14999999999999,153.73904761904865 87.33333333333333,150.44952380952316 89.51666666666665,150.44952380952316 91.69999999999999,147.16000000000003 93.88333333333333,143.87047619047686 96.06666666666666,140.58095238095373 98.24999999999999,137.29142857142824 100.43333333333332,137.29142857142824 102.61666666666666,137.29142857142824 104.79999999999998,134.0019047619051 106.98333333333332,134.0019047619051 109.16666666666666,134.0019047619051 111.35,134.0019047619051 113.53333333333332,134.0019047619051 115.71666666666665,130.71238095238198 117.89999999999999,130.71238095238198 120.08333333333333,130.71238095238198 122.26666666666665,130.71238095238198 124.44999999999999,127.42285714285649 126.63333333333333,127.42285714285649 128.81666666666666,127.42285714285649 131,124.13333333333334 133.1833333333333,124.13333333333334 135.36666666666665,124.13333333333334 137.54999999999998,120.84380952381021 139.73333333333332,117.55428571428706 141.91666666666666,117.55428571428706 144.1,117.55428571428706 146.28333333333333,114.26476190476158 148.46666666666664,114.26476190476158 150.64999999999998,114.26476190476158 152.83333333333331,114.26476190476158 155.01666666666665,114.26476190476158 157.2,110.97523809523844 159.38333333333333,110.97523809523844 161.56666666666666,110.97523809523844 163.74999999999997,110.97523809523844 165.9333333333333,110.97523809523844 168.11666666666665,110.97523809523844 170.29999999999998,107.68571428571529 172.48333333333332,107.68571428571529 174.66666666666666,107.68571428571529 176.85,107.68571428571529 179.0333333333333,107.68571428571529 181.21666666666664,107.68571428571529 183.39999999999998,107.68571428571529 185.58333333333331,107.68571428571529 187.76666666666665,107.68571428571529 189.95,107.68571428571529 192.13333333333333,107.68571428571529 194.31666666666663,107.68571428571529 196.49999999999997,107.68571428571529 198.6833333333333,104.39619047618982 200.86666666666665,104.39619047618982 203.04999999999998,104.39619047618982 205.23333333333332,104.39619047618982 207.41666666666666,104.39619047618982 209.59999999999997,104.39619047618982 211.7833333333333,104.39619047618982 213.96666666666664,104.39619047618982 216.14999999999998,101.10666666666668 218.33333333333331,101.10666666666668 220.51666666666665,101.10666666666668 222.7,101.10666666666668 224.88333333333333,101.10666666666668 227.06666666666663,101.10666666666668 229.24999999999997,101.10666666666668 231.4333333333333,101.10666666666668 233.61666666666665,101.10666666666668 235.79999999999998,97.81714285714354 237.98333333333332,97.81714285714354 240.16666666666666,97.81714285714354 242.34999999999997,97.81714285714354 244.5333333333333,97.81714285714354 246.71666666666664,94.52761904762039 248.89999999999998,94.52761904762039 251.08333333333331,94.52761904762039 253.26666666666665,94.52761904762039 255.45,94.52761904762039 257.6333333333333,94.52761904762039 259.81666666666666,94.52761904762039 262,94.52761904762039 264.1833333333333,94.52761904762039 266.3666666666666,94.52761904762039 268.54999999999995,91.23809523809491 270.7333333333333,91.23809523809491 272.91666666666663,91.23809523809491 275.09999999999997,91.23809523809491 277.2833333333333,91.23809523809491 279.46666666666664,91.23809523809491 281.65,91.23809523809491 283.8333333333333,91.23809523809491 286.01666666666665,91.23809523809491 288.2,91.23809523809491 290.3833333333333,91.23809523809491 292.56666666666666,91.23809523809491 294.75,91.23809523809491 296.9333333333333,91.23809523809491 299.1166666666666,91.23809523809491 301.29999999999995,91.23809523809491 303.4833333333333,91.23809523809491 305.66666666666663,91.23809523809491 307.84999999999997,87.94857142857177 310.0333333333333,87.94857142857177 312.21666666666664,87.94857142857177 314.4,87.94857142857177 316.5833333333333,87.94857142857177 318.76666666666665,87.94857142857177 320.95,87.94857142857177 323.1333333333333,87.94857142857177 325.31666666666666,87.94857142857177 327.49999999999994,87.94857142857177 329.6833333333333,87.94857142857177 331.8666666666666,87.94857142857177 334.04999999999995,87.94857142857177 336.2333333333333,87.94857142857177 338.41666666666663,87.94857142857177 340.59999999999997,87.94857142857177 342.7833333333333,87.94857142857177 344.96666666666664,87.94857142857177 347.15,87.94857142857177 349.3333333333333,87.94857142857177 351.51666666666665,87.94857142857177 353.7,87.94857142857177 355.8833333333333,87.94857142857177 358.0666666666666,87.94857142857177 360.24999999999994,87.94857142857177 362.4333333333333,87.94857142857177 364.6166666666666,87.94857142857177 366.79999999999995,87.94857142857177 368.9833333333333,87.94857142857177 371.16666666666663,88.94857142857177 373.34999999999997,88.94857142857177 375.5333333333333,87.94857142857177 377.71666666666664,87.94857142857177 379.9,88.94857142857177 382.0833333333333,87.94857142857177 384.26666666666665,87.94857142857177 386.45,88.94857142857177 388.63333333333327,88.94857142857177 390.8166666666666,87.94857142857177 392.99999999999994,87.94857142857177 395.1833333333333,87.94857142857177 397.3666666666666,87.94857142857177 399.54999999999995,87.94857142857177 401.7333333333333,87.94857142857177 403.91666666666663,87.94857142857177 406.09999999999997,87.94857142857177 408.2833333333333,87.94857142857177 410.46666666666664,88.94857142857177 412.65,88.94857142857177 414.8333333333333,87.94857142857177 417.01666666666665,87.94857142857177 419.19999999999993,88.94857142857177 421.38333333333327,87.94857142857177 423.5666666666666,87.94857142857177 425.74999999999994,87.94857142857177 427.9333333333333,87.94857142857177 430.1166666666666,87.94857142857177 432.29999999999995,87.94857142857177 434.4833333333333,87.94857142857177 436.66666666666663,87.94857142857177 438.84999999999997,87.94857142857177 441.0333333333333,87.94857142857177 443.21666666666664,87.94857142857177 445.4,84.65904761904864 447.5833333333333,84.65904761904864 449.76666666666665,84.65904761904864 451.94999999999993,84.65904761904864 454.13333333333327,84.65904761904864 456.3166666666666,84.65904761904864 458.49999999999994,84.65904761904864 460.6833333333333,84.65904761904864 462.8666666666666,84.65904761904864 465.04999999999995,84.65904761904864 467.2333333333333,84.65904761904864 469.41666666666663,84.65904761904864 471.59999999999997,85.65904761904864 473.7833333333333,84.65904761904864 475.96666666666664,84.65904761904864 478.15,84.65904761904864 480.3333333333333,84.65904761904864 482.5166666666666,84.65904761904864 484.69999999999993,84.65904761904864 486.88333333333327,84.65904761904864 489.0666666666666,84.65904761904864 491.24999999999994,84.65904761904864 493.4333333333333,84.65904761904864 495.6166666666666,84.65904761904864 497.79999999999995,84.65904761904864 499.9833333333333,84.65904761904864 502.16666666666663,84.65904761904864 504.34999999999997,84.65904761904864 506.5333333333333,84.65904761904864 508.71666666666664,84.65904761904864 510.9,84.65904761904864 513.0833333333333,84.65904761904864 515.2666666666667,84.65904761904864 517.4499999999999,84.65904761904864 519.6333333333333,84.65904761904864 521.8166666666666,84.65904761904864 524,84.65904761904864" style="fill: none;stroke: #EB7113;" stroke-dasharray="" stroke-width="1.2"></polyline>
        
        <polyline points="0,87.94857142857177 524,87.94857142857177" style="fill: none;stroke: rgba(235,51,51,0.5);" stroke-dasharray="4 2" stroke-width="1.2"></polyline><line x1="1.0871369294605808" y1="323" x2="1.0871369294605808" y2="319.58483754512633" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="3.2614107883817427" y1="323" x2="3.2614107883817427" y2="277.0332794223827" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="5.435684647302904" y1="323" x2="5.435684647302904" y2="266.24" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="7.609958506224066" y1="323" x2="7.609958506224066" y2="305.4802166064982" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="9.784232365145227" y1="323" x2="9.784232365145227" y2="312.891119133574" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="11.958506224066388" y1="323" x2="11.958506224066388" y2="318.15046931407943" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="14.132780082987551" y1="323" x2="14.132780082987551" y2="315.48664259927796" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="16.307053941908713" y1="323" x2="16.307053941908713" y2="318.73104693140795" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="18.481327800829874" y1="323" x2="18.481327800829874" y2="309.2027436823105" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="20.655601659751035" y1="323" x2="20.655601659751035" y2="321.0192057761733" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="22.829875518672196" y1="323" x2="22.829875518672196" y2="317.43328519855595" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="25.004149377593357" y1="323" x2="25.004149377593357" y2="318.15046931407943" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="27.178423236514522" y1="323" x2="27.178423236514522" y2="316.20382671480144" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="29.352697095435683" y1="323" x2="29.352697095435683" y2="314.76945848375453" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="31.526970954356845" y1="323" x2="31.526970954356845" y2="314.5303971119134" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="33.70124481327801" y1="323" x2="33.70124481327801" y2="302.91884476534295" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="35.87551867219917" y1="323" x2="35.87551867219917" y2="315.7257039711191" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="38.04979253112033" y1="323" x2="38.04979253112033" y2="312.20808664259926" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="40.22406639004149" y1="323" x2="40.22406639004149" y2="318.52613718411556" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="42.398340248962654" y1="323" x2="42.398340248962654" y2="313.94981949458486" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="44.572614107883815" y1="323" x2="44.572614107883815" y2="315.3158844765343" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="46.74688796680498" y1="323" x2="46.74688796680498" y2="313.7107581227437" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="48.92116182572614" y1="323" x2="48.92116182572614" y2="316.4770397111913" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="51.0954356846473" y1="323" x2="51.0954356846473" y2="319.1408664259928" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="53.26970954356847" y1="323" x2="53.26970954356847" y2="317.945559566787" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="55.44398340248963" y1="323" x2="55.44398340248963" y2="295.61039711191336" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="57.61825726141079" y1="323" x2="57.61825726141079" y2="319.34577617328523" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="59.79253112033195" y1="323" x2="59.79253112033195" y2="318.9701083032491" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="61.96680497925311" y1="323" x2="61.96680497925311" y2="319.1750180505415" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="64.14107883817427" y1="323" x2="64.14107883817427" y2="316.85270758122743" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="66.31535269709543" y1="323" x2="66.31535269709543" y2="302.37241877256315" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="68.48962655601659" y1="323" x2="68.48962655601659" y2="318.56028880866427" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="70.66390041493776" y1="323" x2="70.66390041493776" y2="308.5197111913358" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="72.83817427385893" y1="323" x2="72.83817427385893" y2="305.00209386281585" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="75.01244813278008" y1="323" x2="75.01244813278008" y2="317.1259205776173" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="77.18672199170125" y1="323" x2="77.18672199170125" y2="314.63285198555957" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="79.3609958506224" y1="323" x2="79.3609958506224" y2="293.66375451263536" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="81.53526970954357" y1="323" x2="81.53526970954357" y2="307.7342238267148" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="83.70954356846472" y1="323" x2="83.70954356846472" y2="314.6670036101083" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="85.88381742738589" y1="323" x2="85.88381742738589" y2="317.296678700361" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="88.05809128630705" y1="323" x2="88.05809128630705" y2="308.79292418772565" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="90.23236514522821" y1="323" x2="90.23236514522821" y2="305.41191335740075" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="92.40663900414937" y1="323" x2="92.40663900414937" y2="306.26570397111914" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="94.58091286307054" y1="323" x2="94.58091286307054" y2="312.85696750902525" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="96.75518672199169" y1="323" x2="96.75518672199169" y2="291.7205270758123" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="98.92946058091286" y1="323" x2="98.92946058091286" y2="319.31162454873646" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="101.10373443983401" y1="323" x2="101.10373443983401" y2="312.9594223826715" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="103.27800829875518" y1="323" x2="103.27800829875518" y2="315.48664259927796" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="105.45228215767635" y1="323" x2="105.45228215767635" y2="317.1600722021661" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="107.6265560165975" y1="323" x2="107.6265560165975" y2="313.50584837545125" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="109.80082987551867" y1="323" x2="109.80082987551867" y2="322.3169675090253" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="111.97510373443983" y1="323" x2="111.97510373443983" y2="319.6531407942238" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="114.149377593361" y1="323" x2="114.149377593361" y2="320.23098628158846" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="116.32365145228215" y1="323" x2="116.32365145228215" y2="320.1995667870036" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="118.49792531120332" y1="323" x2="118.49792531120332" y2="317.50158844765343" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="120.67219917012447" y1="323" x2="120.67219917012447" y2="319.6189891696751" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="122.84647302904564" y1="323" x2="122.84647302904564" y2="317.0234657039711" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="125.0207468879668" y1="323" x2="125.0207468879668" y2="314.5303971119134" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="127.19502074688796" y1="323" x2="127.19502074688796" y2="318.2529241877256" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="129.3692946058091" y1="323" x2="129.3692946058091" y2="319.0725631768953" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="131.54356846473027" y1="323" x2="131.54356846473027" y2="308.5880144404332" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="133.71784232365144" y1="323" x2="133.71784232365144" y2="318.76519855595666" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="135.89211618257258" y1="323" x2="135.89211618257258" y2="317.296678700361" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="138.06639004149375" y1="323" x2="138.06639004149375" y2="305.7875812274368" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="140.24066390041492" y1="323" x2="140.24066390041492" y2="295.4737906137184" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="142.41493775933608" y1="323" x2="142.41493775933608" y2="319.1750180505415" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="144.58921161825725" y1="323" x2="144.58921161825725" y2="319.5165342960289" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="146.7634854771784" y1="323" x2="146.7634854771784" y2="316.1696750902527" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="148.93775933609956" y1="323" x2="148.93775933609956" y2="317.7748014440433" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="151.11203319502073" y1="323" x2="151.11203319502073" y2="318.73104693140795" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="153.2863070539419" y1="323" x2="153.2863070539419" y2="318.457833935018" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="155.46058091286304" y1="323" x2="155.46058091286304" y2="317.87725631768956" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="157.6348547717842" y1="323" x2="157.6348547717842" y2="315.41833935018053" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="159.80912863070537" y1="323" x2="159.80912863070537" y2="319.31162454873646" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="161.98340248962654" y1="323" x2="161.98340248962654" y2="310.22729241877255" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="164.15767634854768" y1="323" x2="164.15767634854768" y2="310.80787003610106" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="166.33195020746885" y1="323" x2="166.33195020746885" y2="316.1696750902527" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="168.50622406639002" y1="323" x2="168.50622406639002" y2="315.65740072202163" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="170.6804979253112" y1="323" x2="170.6804979253112" y2="315.24758122743685" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="172.85477178423233" y1="323" x2="172.85477178423233" y2="315.89646209386285" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="175.0290456431535" y1="323" x2="175.0290456431535" y2="307.3585559566787" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="177.20331950207466" y1="323" x2="177.20331950207466" y2="309.2368953068592" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="179.37759336099583" y1="323" x2="179.37759336099583" y2="313.54" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="181.551867219917" y1="323" x2="181.551867219917" y2="318.28707581227434" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="183.72614107883814" y1="323" x2="183.72614107883814" y2="318.28707581227434" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="185.9004149377593" y1="323" x2="185.9004149377593" y2="318.18462093862814" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="188.07468879668048" y1="323" x2="188.07468879668048" y2="319.8238989169675" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="190.24896265560164" y1="323" x2="190.24896265560164" y2="318.4236823104693" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="192.42323651452278" y1="323" x2="192.42323651452278" y2="318.04801444043323" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="194.59751037344395" y1="323" x2="194.59751037344395" y2="317.8431046931408" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="196.77178423236512" y1="323" x2="196.77178423236512" y2="313.94981949458486" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="198.9460580912863" y1="323" x2="198.9460580912863" y2="317.87725631768956" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="201.12033195020743" y1="323" x2="201.12033195020743" y2="319.37992779783394" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="203.2946058091286" y1="323" x2="203.2946058091286" y2="317.1600722021661" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="205.46887966804977" y1="323" x2="205.46887966804977" y2="310.29559566787003" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="207.64315352697093" y1="323" x2="207.64315352697093" y2="315.1109747292419" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="209.8174273858921" y1="323" x2="209.8174273858921" y2="318.9701083032491" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="211.99170124481324" y1="323" x2="211.99170124481324" y2="321.29241877256317" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="214.1659751037344" y1="323" x2="214.1659751037344" y2="308.3831046931408" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="216.34024896265558" y1="323" x2="216.34024896265558" y2="319.6531407942238" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="218.51452282157675" y1="323" x2="218.51452282157675" y2="318.4236823104693" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="220.6887966804979" y1="323" x2="220.6887966804979" y2="318.2187725631769" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="222.86307053941906" y1="323" x2="222.86307053941906" y2="316.4770397111913" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="225.03734439834022" y1="323" x2="225.03734439834022" y2="320.50693140794226" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="227.2116182572614" y1="323" x2="227.2116182572614" y2="318.93595667870034" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="229.38589211618253" y1="323" x2="229.38589211618253" y2="311.8665703971119" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="231.5601659751037" y1="323" x2="231.5601659751037" y2="314.18888086642596" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="233.73443983402487" y1="323" x2="233.73443983402487" y2="319.4482310469314" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="235.90871369294604" y1="323" x2="235.90871369294604" y2="316.3745848375451" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="238.08298755186718" y1="323" x2="238.08298755186718" y2="320.541083032491" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="240.25726141078835" y1="323" x2="240.25726141078835" y2="318.3895306859206" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="242.4315352697095" y1="323" x2="242.4315352697095" y2="318.01386281588447" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="244.60580912863068" y1="323" x2="244.60580912863068" y2="303.7043321299639" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="246.78008298755185" y1="323" x2="246.78008298755185" y2="314.80361010830325" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="248.954356846473" y1="323" x2="248.954356846473" y2="318.457833935018" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="251.12863070539416" y1="323" x2="251.12863070539416" y2="319.4482310469314" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="253.30290456431533" y1="323" x2="253.30290456431533" y2="319.106714801444" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="255.4771784232365" y1="323" x2="255.4771784232365" y2="318.18462093862814" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="257.65145228215766" y1="323" x2="257.65145228215766" y2="310.87617328519855" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="259.82572614107886" y1="323" x2="259.82572614107886" y2="320.4727797833935" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="262" y1="323" x2="262" y2="319.96050541516246" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="264.17427385892114" y1="323" x2="264.17427385892114" y2="315.7598555956679" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="266.34854771784234" y1="323" x2="266.34854771784234" y2="320.2678700361011" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="268.5228215767635" y1="323" x2="268.5228215767635" y2="319.7897472924188" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="270.6970954356846" y1="323" x2="270.6970954356846" y2="320.02880866425994" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="272.8713692946058" y1="323" x2="272.8713692946058" y2="319.0042599277978" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="275.04564315352695" y1="323" x2="275.04564315352695" y2="320.404476534296" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="277.21991701244815" y1="323" x2="277.21991701244815" y2="316.8868592057762" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="279.3941908713693" y1="323" x2="279.3941908713693" y2="319.5506859205776" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="281.56846473029043" y1="323" x2="281.56846473029043" y2="319.48238267148014" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="283.7427385892116" y1="323" x2="283.7427385892116" y2="320.4386281588448" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="285.91701244813277" y1="323" x2="285.91701244813277" y2="319.5506859205776" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="288.09128630705396" y1="323" x2="288.09128630705396" y2="314.05227436823105" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="290.2655601659751" y1="323" x2="290.2655601659751" y2="322.4535740072202" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="292.43983402489624" y1="323" x2="292.43983402489624" y2="314.08642599277977" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="294.61410788381744" y1="323" x2="294.61410788381744" y2="316.8185559566787" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="296.7883817427386" y1="323" x2="296.7883817427386" y2="315.65740072202163" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="298.9626556016597" y1="323" x2="298.9626556016597" y2="317.43328519855595" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="301.1369294605809" y1="323" x2="301.1369294605809" y2="318.56028880866427" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="303.31120331950206" y1="323" x2="303.31120331950206" y2="319.85805054151626" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="305.48547717842325" y1="323" x2="305.48547717842325" y2="317.7406498194946" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="307.6597510373444" y1="323" x2="307.6597510373444" y2="320.5752346570397" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="309.83402489626553" y1="323" x2="309.83402489626553" y2="320.64353790613717" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="312.00829875518673" y1="323" x2="312.00829875518673" y2="320.9509025270758" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="314.18257261410787" y1="323" x2="314.18257261410787" y2="321.25826714801445" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="316.35684647302907" y1="323" x2="316.35684647302907" y2="321.29241877256317" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="318.5311203319502" y1="323" x2="318.5311203319502" y2="313.57415162454873" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="320.70539419087135" y1="323" x2="320.70539419087135" y2="320.78014440433213" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="322.87966804979254" y1="323" x2="322.87966804979254" y2="320.91675090252704" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="325.0539419087137" y1="323" x2="325.0539419087137" y2="320.60938628158846" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="327.2282157676348" y1="323" x2="327.2282157676348" y2="318.7993501805054" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="329.402489626556" y1="323" x2="329.402489626556" y2="312.72036101083035" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="331.57676348547716" y1="323" x2="331.57676348547716" y2="319.85805054151626" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="333.75103734439836" y1="323" x2="333.75103734439836" y2="319.2774729241877" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="335.9253112033195" y1="323" x2="335.9253112033195" y2="319.6531407942238" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="338.09958506224064" y1="323" x2="338.09958506224064" y2="318.8676534296029" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="340.27385892116183" y1="323" x2="340.27385892116183" y2="320.71184115523465" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="342.448132780083" y1="323" x2="342.448132780083" y2="319.1408664259928" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="344.6224066390041" y1="323" x2="344.6224066390041" y2="316.75025270758124" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="346.7966804979253" y1="323" x2="346.7966804979253" y2="319.75559566787" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="348.97095435684645" y1="323" x2="348.97095435684645" y2="320.91675090252704" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="351.14522821576764" y1="323" x2="351.14522821576764" y2="319.5506859205776" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="353.3195020746888" y1="323" x2="353.3195020746888" y2="320.8484476534296" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="355.4937759336099" y1="323" x2="355.4937759336099" y2="320.541083032491" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="357.6680497925311" y1="323" x2="357.6680497925311" y2="319.58483754512633" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="359.84232365145226" y1="323" x2="359.84232365145226" y2="319.0384115523466" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="362.01659751037346" y1="323" x2="362.01659751037346" y2="319.96050541516246" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="364.1908713692946" y1="323" x2="364.1908713692946" y2="321.4973285198556" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="366.36514522821574" y1="323" x2="366.36514522821574" y2="320.74599277978336" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="368.53941908713693" y1="323" x2="368.53941908713693" y2="318.01386281588447" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="370.7136929460581" y1="323" x2="370.7136929460581" y2="311.32014440433215" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="372.8879668049792" y1="323" x2="372.8879668049792" y2="318.2187725631769" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="375.0622406639004" y1="323" x2="375.0622406639004" y2="316.9893140794224" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="377.23651452282155" y1="323" x2="377.23651452282155" y2="317.91140794223827" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="379.41078838174275" y1="323" x2="379.41078838174275" y2="321.7363898916967" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="381.5850622406639" y1="323" x2="381.5850622406639" y2="314.2571841155235" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="383.759336099585" y1="323" x2="383.759336099585" y2="318.8676534296029" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="385.9336099585062" y1="323" x2="385.9336099585062" y2="319.8238989169675" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="388.10788381742736" y1="323" x2="388.10788381742736" y2="316.57949458483756" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="390.28215767634856" y1="323" x2="390.28215767634856" y2="320.88259927797833" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="392.4564315352697" y1="323" x2="392.4564315352697" y2="320.60938628158846" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="394.63070539419084" y1="323" x2="394.63070539419084" y2="316.4428880866426" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="396.80497925311204" y1="323" x2="396.80497925311204" y2="314.80361010830325" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="398.9792531120332" y1="323" x2="398.9792531120332" y2="318.93595667870034" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="401.1535269709543" y1="323" x2="401.1535269709543" y2="314.9743682310469" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="403.3278008298755" y1="323" x2="403.3278008298755" y2="320.4386281588448" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="405.50207468879665" y1="323" x2="405.50207468879665" y2="320.64353790613717" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="407.67634854771785" y1="323" x2="407.67634854771785" y2="320.88259927797833" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="409.850622406639" y1="323" x2="409.850622406639" y2="318.28707581227434" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="412.02489626556013" y1="323" x2="412.02489626556013" y2="318.69655379061373" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="414.1991701244813" y1="323" x2="414.1991701244813" y2="319.48238267148014" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="416.37344398340247" y1="323" x2="416.37344398340247" y2="317.6040433212996" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="418.54771784232366" y1="323" x2="418.54771784232366" y2="319.96050541516246" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="420.7219917012448" y1="323" x2="420.7219917012448" y2="278.7053429602888" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="422.89626556016594" y1="323" x2="422.89626556016594" y2="317.6723465703971" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="425.07053941908714" y1="323" x2="425.07053941908714" y2="314.63285198555957" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="427.2448132780083" y1="323" x2="427.2448132780083" y2="318.2529241877256" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="429.4190871369294" y1="323" x2="429.4190871369294" y2="319.5506859205776" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="431.5933609958506" y1="323" x2="431.5933609958506" y2="317.6040433212996" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="433.76763485477176" y1="323" x2="433.76763485477176" y2="317.1259205776173" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="435.94190871369295" y1="323" x2="435.94190871369295" y2="320.1654151624549" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="438.1161825726141" y1="323" x2="438.1161825726141" y2="314.6670036101083" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="440.29045643153523" y1="323" x2="440.29045643153523" y2="319.96050541516246" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="442.46473029045643" y1="323" x2="442.46473029045643" y2="318.2529241877256" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="444.63900414937757" y1="323" x2="444.63900414937757" y2="317.5698916967509" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="446.8132780082987" y1="323" x2="446.8132780082987" y2="320.8484476534296" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="448.9875518672199" y1="323" x2="448.9875518672199" y2="319.4482310469314" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="451.16182572614105" y1="323" x2="451.16182572614105" y2="318.11631768953066" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="453.33609958506224" y1="323" x2="453.33609958506224" y2="317.7406498194946" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="455.5103734439834" y1="323" x2="455.5103734439834" y2="320.13126353790614" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="457.6846473029045" y1="323" x2="457.6846473029045" y2="313.7107581227437" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="459.8589211618257" y1="323" x2="459.8589211618257" y2="318.457833935018" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="462.03319502074686" y1="323" x2="462.03319502074686" y2="316.20382671480144" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="464.20746887966806" y1="323" x2="464.20746887966806" y2="312.5837545126354" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="466.3817427385892" y1="323" x2="466.3817427385892" y2="313.54" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="468.55601659751034" y1="323" x2="468.55601659751034" y2="320.09711191335737" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="470.73029045643153" y1="323" x2="470.73029045643153" y2="318.18462093862814" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="472.9045643153527" y1="323" x2="472.9045643153527" y2="316.8868592057762" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="475.0788381742738" y1="323" x2="475.0788381742738" y2="318.56028880866427" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="477.253112033195" y1="323" x2="477.253112033195" y2="314.80361010830325" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="479.42738589211615" y1="323" x2="479.42738589211615" y2="319.2774729241877" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="481.60165975103735" y1="323" x2="481.60165975103735" y2="317.7748014440433" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="483.7759336099585" y1="323" x2="483.7759336099585" y2="319.6531407942238" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="485.9502074688796" y1="323" x2="485.9502074688796" y2="318.76519855595666" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="488.1244813278008" y1="323" x2="488.1244813278008" y2="320.71184115523465" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="490.29875518672196" y1="323" x2="490.29875518672196" y2="319.96050541516246" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="492.47302904564316" y1="323" x2="492.47302904564316" y2="312.82281588447654" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="494.6473029045643" y1="323" x2="494.6473029045643" y2="320.404476534296" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="496.82157676348544" y1="323" x2="496.82157676348544" y2="312.3105415162455" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="498.99585062240664" y1="323" x2="498.99585062240664" y2="316.68194945848376" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="501.1701244813278" y1="323" x2="501.1701244813278" y2="302.7480866425993" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="503.3443983402489" y1="323" x2="503.3443983402489" y2="308.0074368231047" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="505.5186721991701" y1="323" x2="505.5186721991701" y2="319.0042599277978" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="507.69294605809125" y1="323" x2="507.69294605809125" y2="318.2761472924188" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="509.86721991701245" y1="323" x2="509.86721991701245" y2="317.43328519855595" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="512.0414937759335" y1="323" x2="512.0414937759335" y2="315.38418772563176" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="514.2157676348547" y1="323" x2="514.2157676348547" y2="317.945559566787" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="516.3900414937759" y1="323" x2="516.3900414937759" y2="308.7171075812274" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="518.5643153526971" y1="323" x2="518.5643153526971" y2="322.21451263537904" style="stroke: #D60A22; stroke-width: 2.1742738589211617px;"></line><line x1="520.7385892116182" y1="323" x2="520.7385892116182" y2="323" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><line x1="522.9128630705394" y1="323" x2="522.9128630705394" y2="296.43003610108303" style="stroke: #037B66; stroke-width: 2.1742738589211617px;"></line><text x="0" y="255.24" fill="#3D404D" stroke="#fff" stroke-width="0" text-anchor="start" writing-mode="" letter-spacing="0" dominant-baseline="text-before-edge" paint-order="stroke" style="font-size: 12px; font-weight: 400; font-family: ;">成交量778手</text><text x="0" y="3" fill="#B7B9C1" dominant-baseline="text-before-edge" text-anchor="start" style="font-size: 12px; font-weight: 400; font-family: ;">36.38</text>,<text x="0" y="112.62" fill="#B7B9C1" dominant-baseline="middle" text-anchor="start" style="font-size: 12px; font-weight: 400; font-family: ;">36.06</text>,<text x="0" y="222.24" fill="#B7B9C1" dominant-baseline="" text-anchor="start" style="font-size: 12px; font-weight: 400; font-family: ;">35.75</text><text x="524" y="3" fill="#B7B9C1" dominant-baseline="text-before-edge" text-anchor="end" style="font-size: 12px; font-weight: 400; font-family: ;">1.51%</text>,<text x="524" y="112.62" fill="#B7B9C1" dominant-baseline="middle" text-anchor="end" style="font-size: 12px; font-weight: 400; font-family: ;">0.63%</text>,<text x="524" y="222.24" fill="#B7B9C1" dominant-baseline="" text-anchor="end" style="font-size: 12px; font-weight: 400; font-family: ;">-0.25%</text></svg><!--887--></div><!--370-->
                                <!--371-->
                                <div class="feedback" data-s-9dda213a=""></div>

                                <div class="feedback-svg" data-s-9dda213a=""></div>
                                <!--372-->
                                <!--373-->
                                
                                <!--374-->
                                <!--375-->
                                <!--376-->
                            </div><!--369-->
                        </div>
                    <!--368-->
    <object tabindex="-1" type="text/html" aria-hidden="true" data="about:blank" style="display: block; position: absolute; top: 0px; left: 0px; width: 100%; height: 100%; border: none; padding: 0px; margin: 0px; opacity: 0; z-index: -1000; pointer-events: none;"></object></div>
                </div>

                <div class="chart-control-bar-container" data-s-9dda213a="">
                    <div data-s-0bd93592="" class="match-module-container" id="match-module-container-fac-market-chart-new-chart-thumbnail-0.10445714114903759" style="position: relative;">
        <!--379-->
                        
                        <div class="chart-control-bar" data-s-9dda213a=""><!--380--><svg version="1.1" xmlns="http://www.w3.org/2000/svg" width="524px" height="47px"><line x1="1" y1="1" x2="524" y2="1" stroke-dasharray="" style="stroke: #fff; stroke-width: 1px;"></line>,<line x1="1" y1="29" x2="524" y2="29" stroke-dasharray="" style="stroke: #fff; stroke-width: 1px;"></line>,<line x1="0" y1="0" x2="0" y2="29" stroke-dasharray="" style="stroke: #fff; stroke-width: 1px;"></line>,<line x1="104.8" y1="0" x2="104.8" y2="29" stroke-dasharray="" style="stroke: #fff; stroke-width: 1px;"></line>,<line x1="209.6" y1="0" x2="209.6" y2="29" stroke-dasharray="" style="stroke: #fff; stroke-width: 1px;"></line>,<line x1="314.4" y1="0" x2="314.4" y2="29" stroke-dasharray="" style="stroke: #fff; stroke-width: 1px;"></line>,<line x1="419.2" y1="0" x2="419.2" y2="29" stroke-dasharray="" style="stroke: #fff; stroke-width: 1px;"></line><text x="52.4" y="45" fill="#848691" text-anchor="middle" style="font-size: 12px; font-weight: 400; font-family: ;">04-24</text>,<text x="157.2" y="45" fill="#848691" text-anchor="middle" style="font-size: 12px; font-weight: 400; font-family: ;">04-25</text>,<text x="262" y="45" fill="#848691" text-anchor="middle" style="font-size: 12px; font-weight: 400; font-family: ;">04-28</text>,<text x="366.79999999999995" y="45" fill="#848691" text-anchor="middle" style="font-size: 12px; font-weight: 400; font-family: ;">04-29</text>,<text x="471.59999999999997" y="45" fill="#848691" text-anchor="middle" style="font-size: 12px; font-weight: 400; font-family: ;">04-30</text>
        <defs>
                <linearGradient id="market_chart_bar" x1="0%" y1="0%" x2="0%" y2="100%">
                    <stop offset="0%" stop-color="rgb(78, 110, 242, .2)"></stop>
                    <stop offset="100%" stop-color="rgba(255, 255, 255, 0)" stop-opacity="0"></stop>
                </linearGradient>
            </defs>
            <polyline points="0,22.82978723404245 0.43666666666666665,25.914893617021225 0.8733333333333333,23.755319148936103 1.31,22.521276595744638 1.7466666666666666,22.212765957446827 2.183333333333333,21.904255319149016 2.62,21.595744680850984 3.0566666666666666,20.670212765957547 3.493333333333333,19.436170212765862 3.9299999999999997,19.436170212765862 4.366666666666666,20.670212765957547 4.803333333333333,21.287234042553173 5.24,20.97872340425536 5.676666666666667,21.904255319149016 6.113333333333333,22.521276595744638 6.55,22.521276595744638 6.986666666666666,21.287234042553173 7.423333333333333,20.053191489361705 7.859999999999999,19.12765957446805 8.296666666666667,20.97872340425536 8.733333333333333,20.97872340425536 9.17,20.670212765957547 9.606666666666666,21.904255319149016 10.043333333333333,21.595744680850984 10.48,19.436170212765862 10.916666666666666,20.053191489361705 11.353333333333333,19.12765957446805 11.79,16.968085106382933 12.226666666666667,15.734042553191465 12.663333333333332,15.425531914893654 13.1,16.351063829787307 13.536666666666665,15.734042553191465 13.973333333333333,16.042553191489276 14.41,15.734042553191465 14.846666666666666,15.734042553191465 15.283333333333333,17.276595744680744 15.719999999999999,17.893617021276587 16.156666666666666,16.968085106382933 16.593333333333334,16.968085106382933 17.029999999999998,17.276595744680744 17.466666666666665,16.351063829787307 17.903333333333332,16.65957446808512 18.34,15.425531914893654 18.776666666666667,15.734042553191465 19.21333333333333,15.117021276595622 19.65,14.808510638297811 20.086666666666666,14.808510638297811 20.523333333333333,15.734042553191465 20.96,16.351063829787307 21.396666666666665,17.585106382978775 21.833333333333332,16.65957446808512 22.27,15.425531914893654 22.706666666666667,17.276595744680744 23.14333333333333,16.968085106382933 23.58,16.65957446808512 24.016666666666666,16.65957446808512 24.453333333333333,15.734042553191465 24.89,16.351063829787307 25.326666666666664,16.042553191489276 25.763333333333332,15.425531914893654 26.2,15.425531914893654 26.636666666666667,13.574468085106345 27.07333333333333,14.191489361702187 27.509999999999998,16.65957446808512 27.946666666666665,16.351063829787307 28.383333333333333,17.276595744680744 28.82,16.351063829787307 29.256666666666664,16.042553191489276 29.69333333333333,17.276595744680744 30.13,17.276595744680744 30.566666666666666,16.351063829787307 31.00333333333333,16.65957446808512 31.439999999999998,16.351063829787307 31.876666666666665,15.425531914893654 32.31333333333333,16.042553191489276 32.75,17.585106382978775 33.18666666666667,18.202127659574398 33.623333333333335,16.968085106382933 34.059999999999995,16.65957446808512 34.49666666666666,16.65957446808512 34.93333333333333,17.276595744680744 35.37,17.893617021276587 35.806666666666665,18.81914893617024 36.24333333333333,17.276595744680744 36.68,18.202127659574398 37.11666666666667,19.436170212765862 37.553333333333335,19.744680851063894 37.989999999999995,19.436170212765862 38.42666666666666,20.670212765957547 38.86333333333333,20.053191489361705 39.3,19.744680851063894 39.736666666666665,21.904255319149016 40.17333333333333,18.51063829787243 40.61,19.12765957446805 41.04666666666667,20.053191489361705 41.483333333333334,21.287234042553173 41.92,22.82978723404245 42.35666666666666,23.13829787234048 42.79333333333333,24.063829787234134 43.23,25.297872340425602 43.666666666666664,23.13829787234048 44.10333333333333,23.44680851063829 44.54,24.063829787234134 44.97666666666667,24.372340425531945 45.413333333333334,23.755319148936103 45.85,24.372340425531945 46.28666666666666,25.606382978723413 46.72333333333333,27.14893617021269 47.16,24.68085106382976 47.596666666666664,24.063829787234134 48.03333333333333,23.755319148936103 48.47,23.44680851063829 48.906666666666666,23.13829787234048 49.343333333333334,23.13829787234048 49.78,22.82978723404245 50.21666666666666,22.212765957446827 50.65333333333333,21.595744680850984 51.089999999999996,22.521276595744638 51.526666666666664,22.521276595744638 51.96333333333333,22.521276595744638 52.4,23.44680851063829 52.836666666666666,21.595744680850984 53.27333333333333,23.13829787234048 53.71,23.755319148936103 54.14666666666666,24.98936170212757 54.58333333333333,24.372340425531945 55.019999999999996,23.755319148936103 55.45666666666666,23.44680851063829 55.89333333333333,23.44680851063829 56.33,22.82978723404245 56.766666666666666,22.82978723404245 57.20333333333333,23.13829787234048 57.64,23.44680851063829 58.07666666666666,23.755319148936103 58.51333333333333,22.521276595744638 58.949999999999996,21.287234042553173 59.38666666666666,21.595744680850984 59.82333333333333,23.755319148936103 60.26,23.44680851063829 60.696666666666665,22.521276595744638 61.13333333333333,22.212765957446827 61.57,22.82978723404245 62.00666666666666,21.904255319149016 62.44333333333333,21.595744680850984 62.879999999999995,21.904255319149016 63.31666666666666,21.595744680850984 63.75333333333333,21.595744680850984 64.19,20.97872340425536 64.62666666666667,20.97872340425536 65.06333333333333,20.97872340425536 65.5,20.053191489361705 65.93666666666667,20.361702127659516 66.37333333333333,19.744680851063894 66.81,19.436170212765862 67.24666666666667,19.436170212765862 67.68333333333334,19.12765957446805 68.11999999999999,19.12765957446805 68.55666666666666,19.744680851063894 68.99333333333333,17.276595744680744 69.42999999999999,16.042553191489276 69.86666666666666,15.425531914893654 70.30333333333333,15.734042553191465 70.74,16.351063829787307 71.17666666666666,15.734042553191465 71.61333333333333,15.425531914893654 72.05,14.808510638297811 72.48666666666666,14.5 72.92333333333333,13.265957446808534 73.36,12.648936170212691 73.79666666666667,14.191489361702187 74.23333333333333,16.042553191489276 74.67,15.425531914893654 75.10666666666667,15.425531914893654 75.54333333333334,14.808510638297811 75.97999999999999,14.5 76.41666666666666,14.808510638297811 76.85333333333332,15.117021276595622 77.28999999999999,15.425531914893654 77.72666666666666,15.117021276595622 78.16333333333333,15.734042553191465 78.6,15.734042553191465 79.03666666666666,15.425531914893654 79.47333333333333,15.425531914893654 79.91,15.425531914893654 80.34666666666666,14.191489361702187 80.78333333333333,15.117021276595622 81.22,12.648936170212691 81.65666666666667,12.031914893617067 82.09333333333333,10.489361702127571 82.53,11.106382978723413 82.96666666666667,12.34042553191488 83.40333333333334,10.7978723404256 83.84,10.489361702127571 84.27666666666666,9.872340425531947 84.71333333333332,10.180851063829758 85.14999999999999,11.106382978723413 85.58666666666666,11.106382978723413 86.02333333333333,11.414893617021225 86.46,10.7978723404256 86.89666666666666,10.180851063829758 87.33333333333333,11.723404255319036 87.77,13.265957446808534 88.20666666666666,12.648936170212691 88.64333333333333,13.265957446808534 89.08,13.574468085106345 89.51666666666667,14.5 89.95333333333333,14.808510638297811 90.39,15.425531914893654 90.82666666666667,13.574468085106345 91.26333333333334,14.5 91.7,14.191489361702187 92.13666666666666,14.808510638297811 92.57333333333332,14.808510638297811 93.00999999999999,14.191489361702187 93.44666666666666,14.191489361702187 93.88333333333333,12.957446808510722 94.32,13.882978723404158 94.75666666666666,12.34042553191488 95.19333333333333,12.648936170212691 95.63,12.957446808510722 96.06666666666666,12.957446808510722 96.50333333333333,12.34042553191488 96.94,11.106382978723413 97.37666666666667,10.489361702127571 97.81333333333333,10.7978723404256 98.25,11.723404255319036 98.68666666666667,13.265957446808534 99.12333333333333,12.031914893617067 99.56,11.723404255319036 99.99666666666666,12.648936170212691 100.43333333333332,12.648936170212691 100.86999999999999,12.648936170212691 101.30666666666666,12.031914893617067 101.74333333333333,12.031914893617067 102.17999999999999,11.106382978723413 102.61666666666666,10.7978723404256 103.05333333333333,11.106382978723413 103.49,11.106382978723413 103.92666666666666,11.106382978723413 104.36333333333333,10.180851063829758 104.8,9.872340425531947 105.23666666666666,11.414893617021225 105.67333333333333,11.106382978723413 106.11,11.106382978723413 106.54666666666667,8.63829787234048 106.98333333333333,7.095744680850984 107.42,8.021276595744638 107.85666666666665,11.414893617021225 108.29333333333332,13.574468085106345 108.72999999999999,12.957446808510722 109.16666666666666,14.5 109.60333333333332,13.882978723404158 110.03999999999999,14.5 110.47666666666666,12.957446808510722 110.91333333333333,13.265957446808534 111.35,11.414893617021225 111.78666666666666,13.265957446808534 112.22333333333333,12.957446808510722 112.66,12.031914893617067 113.09666666666666,11.106382978723413 113.53333333333333,12.957446808510722 113.97,10.7978723404256 114.40666666666667,9.872340425531947 114.84333333333333,6.787234042553172 115.28,5.244680851063895 115.71666666666667,5.553191489361707 116.15333333333332,4.936170212765864 116.58999999999999,2.468085106382932 117.02666666666666,4.01063829787221 117.46333333333332,6.478723404255361 117.89999999999999,5.244680851063895 118.33666666666666,4.31914893617024 118.77333333333333,5.244680851063895 119.21,3.702127659574398 119.64666666666666,2.15957446808512 120.08333333333333,4.01063829787221 120.52,3.393617021276586 120.95666666666666,0 121.39333333333333,0 121.83,1.5425531914892776 122.26666666666667,2.15957446808512 122.70333333333333,2.15957446808512 123.14,2.15957446808512 123.57666666666667,2.7765957446807437 124.01333333333332,2.15957446808512 124.44999999999999,4.31914893617024 124.88666666666666,3.0851063829787746 125.32333333333332,3.702127659574398 125.75999999999999,4.627659574468052 126.19666666666666,4.627659574468052 126.63333333333333,5.861702127659518 127.07,6.787234042553172 127.50666666666666,8.021276595744638 127.94333333333333,8.32978723404245 128.38,8.946808510638293 128.81666666666666,10.180851063829758 129.25333333333333,10.7978723404256 129.69,9.255319148936104 130.12666666666667,10.180851063829758 130.56333333333333,8.021276595744638 131,8.946808510638293 131.43666666666667,9.255319148936104 131.87333333333333,9.563829787233916 132.31,10.180851063829758 132.74666666666667,10.180851063829758 133.18333333333334,10.7978723404256 133.62,10.489361702127571 134.05666666666667,11.106382978723413 134.49333333333334,11.723404255319036 134.93,9.563829787233916 135.36666666666667,9.563829787233916 135.8033333333333,10.180851063829758 136.23999999999998,10.489361702127571 136.67666666666665,10.180851063829758 137.11333333333332,9.872340425531947 137.54999999999998,9.255319148936104 137.98666666666665,9.872340425531947 138.42333333333332,8.63829787234048 138.85999999999999,8.63829787234048 139.29666666666665,9.563829787233916 139.73333333333332,9.563829787233916 140.17,9.255319148936104 140.60666666666665,9.872340425531947 141.04333333333332,10.489361702127571 141.48,10.7978723404256 141.91666666666666,10.7978723404256 142.35333333333332,10.489361702127571 142.79,11.106382978723413 143.22666666666666,10.7978723404256 143.66333333333333,11.723404255319036 144.1,12.957446808510722 144.53666666666666,12.34042553191488 144.97333333333333,12.031914893617067 145.41,12.031914893617067 145.84666666666666,12.031914893617067 146.28333333333333,10.489361702127571 146.72,11.723404255319036 147.15666666666667,12.031914893617067 147.59333333333333,11.723404255319036 148.03,12.34042553191488 148.46666666666667,11.106382978723413 148.90333333333334,9.872340425531947 149.34,9.563829787233916 149.77666666666667,9.563829787233916 150.21333333333334,9.563829787233916 150.65,10.180851063829758 151.08666666666667,12.34042553191488 151.52333333333334,12.957446808510722 151.95999999999998,14.191489361702187 152.39666666666665,13.574468085106345 152.83333333333331,12.031914893617067 153.26999999999998,12.648936170212691 153.70666666666665,11.723404255319036 154.14333333333332,12.957446808510722 154.57999999999998,12.957446808510722 155.01666666666665,13.882978723404158 155.45333333333332,13.882978723404158 155.89,13.574468085106345 156.32666666666665,13.574468085106345 156.76333333333332,13.882978723404158 157.2,16.968085106382933 157.63666666666666,16.65957446808512 158.07333333333332,15.734042553191465 158.51,15.734042553191465 158.94666666666666,16.042553191489276 159.38333333333333,16.042553191489276 159.82,15.425531914893654 160.25666666666666,15.425531914893654 160.69333333333333,16.042553191489276 161.13,16.042553191489276 161.56666666666666,16.042553191489276 162.00333333333333,14.5 162.44,15.117021276595622 162.87666666666667,13.882978723404158 163.31333333333333,14.808510638297811 163.75,13.574468085106345 164.18666666666667,13.882978723404158 164.62333333333333,13.882978723404158 165.06,13.882978723404158 165.49666666666667,14.191489361702187 165.93333333333334,12.648936170212691 166.37,10.489361702127571 166.80666666666667,11.723404255319036 167.24333333333334,11.414893617021225 167.68,10.489361702127571 168.11666666666665,10.180851063829758 168.5533333333333,10.7978723404256 168.98999999999998,10.180851063829758 169.42666666666665,10.489361702127571 169.86333333333332,9.872340425531947 170.29999999999998,12.648936170212691 170.73666666666665,12.34042553191488 171.17333333333332,12.648936170212691 171.60999999999999,10.489361702127571 172.04666666666665,8.63829787234048 172.48333333333332,9.563829787233916 172.92,8.021276595744638 173.35666666666665,8.021276595744638 173.79333333333332,8.021276595744638 174.23,5.861702127659518 174.66666666666666,5.861702127659518 175.10333333333332,6.787234042553172 175.54,4.936170212765864 175.97666666666666,8.021276595744638 176.41333333333333,8.63829787234048 176.85,8.32978723404245 177.28666666666666,8.021276595744638 177.72333333333333,8.946808510638293 178.16,10.489361702127571 178.59666666666666,8.946808510638293 179.03333333333333,8.63829787234048 179.47,8.021276595744638 179.90666666666667,8.946808510638293 180.34333333333333,8.63829787234048 180.78,9.563829787233916 181.21666666666667,9.255319148936104 181.65333333333334,10.489361702127571 182.09,11.106382978723413 182.52666666666667,10.7978723404256 182.96333333333334,11.106382978723413 183.4,11.106382978723413 183.83666666666664,10.489361702127571 184.2733333333333,11.414893617021225 184.70999999999998,11.106382978723413 185.14666666666665,11.414893617021225 185.58333333333331,11.414893617021225 186.01999999999998,11.723404255319036 186.45666666666665,12.031914893617067 186.89333333333332,13.882978723404158 187.32999999999998,14.191489361702187 187.76666666666665,14.5 188.20333333333332,14.808510638297811 188.64,14.191489361702187 189.07666666666665,14.5 189.51333333333332,14.191489361702187 189.95,14.5 190.38666666666666,14.5 190.82333333333332,14.808510638297811 191.26,13.882978723404158 191.69666666666666,14.808510638297811 192.13333333333333,14.808510638297811 192.57,14.191489361702187 193.00666666666666,14.808510638297811 193.44333333333333,14.5 193.88,14.808510638297811 194.31666666666666,14.808510638297811 194.75333333333333,15.117021276595622 195.19,12.031914893617067 195.62666666666667,12.031914893617067 196.06333333333333,11.106382978723413 196.5,11.723404255319036 196.93666666666667,11.414893617021225 197.37333333333333,12.34042553191488 197.81,13.265957446808534 198.24666666666667,11.414893617021225 198.68333333333334,10.7978723404256 199.12,11.723404255319036 199.55666666666667,13.574468085106345 199.9933333333333,12.031914893617067 200.42999999999998,12.957446808510722 200.86666666666665,12.031914893617067 201.3033333333333,12.957446808510722 201.73999999999998,15.117021276595622 202.17666666666665,16.65957446808512 202.61333333333332,17.585106382978775 203.04999999999998,16.042553191489276 203.48666666666665,17.276595744680744 203.92333333333332,16.351063829787307 204.35999999999999,17.276595744680744 204.79666666666665,18.202127659574398 205.23333333333332,18.51063829787243 205.67,17.893617021276587 206.10666666666665,16.968085106382933 206.54333333333332,17.585106382978775 206.98,17.585106382978775 207.41666666666666,19.436170212765862 207.85333333333332,19.436170212765862 208.29,19.436170212765862 208.72666666666666,19.436170212765862 209.16333333333333,17.893617021276587 209.6,13.574468085106345 210.03666666666666,14.5 210.47333333333333,14.5 210.91,14.5 211.34666666666666,17.893617021276587 211.78333333333333,19.12765957446805 212.22,18.81914893617024 212.65666666666667,19.436170212765862 213.09333333333333,17.893617021276587 213.53,17.276595744680744 213.96666666666667,16.968085106382933 214.40333333333334,18.51063829787243 214.84,18.81914893617024 215.27666666666667,18.202127659574398 215.7133333333333,17.276595744680744 216.14999999999998,16.351063829787307 216.58666666666664,17.893617021276587 217.0233333333333,16.351063829787307 217.45999999999998,15.425531914893654 217.89666666666665,16.351063829787307 218.33333333333331,16.968085106382933 218.76999999999998,15.117021276595622 219.20666666666665,15.117021276595622 219.64333333333332,15.117021276595622 220.07999999999998,16.65957446808512 220.51666666666665,17.893617021276587 220.95333333333332,18.202127659574398 221.39,15.734042553191465 221.82666666666665,17.585106382978775 222.26333333333332,16.351063829787307 222.7,16.968085106382933 223.13666666666666,16.968085106382933 223.57333333333332,19.744680851063894 224.01,20.053191489361705 224.44666666666666,20.053191489361705 224.88333333333333,17.893617021276587 225.32,17.276595744680744 225.75666666666666,17.276595744680744 226.19333333333333,17.276595744680744 226.63,18.202127659574398 227.06666666666666,18.51063829787243 227.50333333333333,19.744680851063894 227.94,19.744680851063894 228.37666666666667,19.744680851063894 228.81333333333333,19.436170212765862 229.25,19.744680851063894 229.68666666666667,19.436170212765862 230.12333333333333,19.12765957446805 230.56,18.81914893617024 230.99666666666667,18.202127659574398 231.43333333333334,17.585106382978775 231.86999999999998,16.968085106382933 232.30666666666664,18.51063829787243 232.7433333333333,18.51063829787243 233.17999999999998,18.202127659574398 233.61666666666665,17.276595744680744 234.0533333333333,16.65957446808512 234.48999999999998,16.042553191489276 234.92666666666665,14.5 235.36333333333332,14.5 235.79999999999998,14.808510638297811 236.23666666666665,15.117021276595622 236.67333333333332,15.425531914893654 237.10999999999999,16.042553191489276 237.54666666666665,16.042553191489276 237.98333333333332,13.882978723404158 238.42,14.191489361702187 238.85666666666665,13.574468085106345 239.29333333333332,13.574468085106345 239.73,11.106382978723413 240.16666666666666,10.489361702127571 240.60333333333332,10.489361702127571 241.04,10.180851063829758 241.47666666666666,8.63829787234048 241.91333333333333,8.946808510638293 242.35,10.180851063829758 242.78666666666666,10.180851063829758 243.22333333333333,11.106382978723413 243.66,10.7978723404256 244.09666666666666,10.180851063829758 244.53333333333333,10.489361702127571 244.97,8.021276595744638 245.40666666666667,5.553191489361707 245.84333333333333,6.787234042553172 246.28,6.17021276595733 246.71666666666667,5.861702127659518 247.15333333333334,6.787234042553172 247.58999999999997,7.095744680850984 248.02666666666664,7.712765957446827 248.4633333333333,8.021276595744638 248.89999999999998,7.095744680850984 249.33666666666664,6.787234042553172 249.7733333333333,7.095744680850984 250.20999999999998,7.095744680850984 250.64666666666665,6.478723404255361 251.08333333333331,6.478723404255361 251.51999999999998,5.553191489361707 251.95666666666665,5.861702127659518 252.39333333333332,6.17021276595733 252.82999999999998,6.478723404255361 253.26666666666665,6.17021276595733 253.70333333333332,5.861702127659518 254.14,5.861702127659518 254.57666666666665,4.01063829787221 255.01333333333332,4.936170212765864 255.45,5.244680851063895 255.88666666666666,4.936170212765864 256.3233333333333,4.936170212765864 256.76,5.553191489361707 257.19666666666666,3.0851063829787746 257.6333333333333,3.702127659574398 258.07,3.393617021276586 258.50666666666666,4.01063829787221 258.9433333333333,4.627659574468052 259.38,4.627659574468052 259.81666666666666,4.31914893617024 260.25333333333333,4.936170212765864 260.69,4.936170212765864 261.12666666666667,5.244680851063895 261.56333333333333,2.7765957446807437 262,0.9255319148936543 262.43666666666667,2.468085106382932 262.87333333333333,1.234042553191466 263.31,0.30851063829781167 263.74666666666667,1.5425531914892776 264.18333333333334,1.5425531914892776 264.62,1.8510638297870894 265.0566666666667,2.15957446808512 265.49333333333334,2.7765957446807437 265.93,3.702127659574398 266.3666666666667,3.393617021276586 266.80333333333334,4.01063829787221 267.24,4.01063829787221 267.6766666666667,4.936170212765864 268.11333333333334,3.0851063829787746 268.55,3.702127659574398 268.9866666666667,4.01063829787221 269.42333333333335,3.702127659574398 269.86,3.393617021276586 270.2966666666667,4.01063829787221 270.73333333333335,4.936170212765864 271.17,3.702127659574398 271.6066666666666,3.0851063829787746 272.0433333333333,4.01063829787221 272.47999999999996,3.702127659574398 272.91666666666663,4.01063829787221 273.3533333333333,4.01063829787221 273.78999999999996,4.01063829787221 274.22666666666663,4.31914893617024 274.6633333333333,4.31914893617024 275.09999999999997,5.244680851063895 275.53666666666663,5.244680851063895 275.9733333333333,5.861702127659518 276.40999999999997,5.861702127659518 276.84666666666664,5.861702127659518 277.2833333333333,6.478723404255361 277.71999999999997,6.478723404255361 278.15666666666664,6.17021276595733 278.5933333333333,6.17021276595733 279.03,6.478723404255361 279.46666666666664,6.17021276595733 279.9033333333333,6.787234042553172 280.34,6.478723404255361 280.77666666666664,7.095744680850984 281.2133333333333,7.095744680850984 281.65,7.404255319148796 282.08666666666664,8.021276595744638 282.5233333333333,8.021276595744638 282.96,7.712765957446827 283.39666666666665,7.404255319148796 283.8333333333333,7.095744680850984 284.27,7.095744680850984 284.70666666666665,5.553191489361707 285.1433333333333,6.17021276595733 285.58,5.861702127659518 286.01666666666665,6.17021276595733 286.4533333333333,6.17021276595733 286.89,6.478723404255361 287.32666666666665,6.787234042553172 287.7633333333333,6.787234042553172 288.2,6.787234042553172 288.63666666666666,6.787234042553172 289.0733333333333,6.478723404255361 289.51,6.478723404255361 289.94666666666666,6.478723404255361 290.3833333333333,7.095744680850984 290.82,7.095744680850984 291.25666666666666,6.478723404255361 291.6933333333333,6.478723404255361 292.13,6.787234042553172 292.56666666666666,6.787234042553172 293.00333333333333,7.095744680850984 293.44,7.404255319148796 293.87666666666667,7.712765957446827 294.31333333333333,8.021276595744638 294.75,8.021276595744638 295.18666666666667,7.712765957446827 295.62333333333333,7.712765957446827 296.06,7.404255319148796 296.49666666666667,7.404255319148796 296.93333333333334,8.32978723404245 297.37,8.32978723404245 297.8066666666667,7.095744680850984 298.24333333333334,7.095744680850984 298.68,7.095744680850984 299.1166666666667,6.478723404255361 299.55333333333334,6.787234042553172 299.99,6.787234042553172 300.4266666666667,6.787234042553172 300.86333333333334,6.787234042553172 301.3,6.478723404255361 301.7366666666667,6.787234042553172 302.17333333333335,6.787234042553172 302.61,6.478723404255361 303.0466666666667,6.478723404255361 303.4833333333333,6.17021276595733 303.91999999999996,3.702127659574398 304.3566666666666,4.01063829787221 304.7933333333333,4.31914893617024 305.22999999999996,4.31914893617024 305.66666666666663,3.393617021276586 306.1033333333333,2.468085106382932 306.53999999999996,4.936170212765864 306.97666666666663,4.936170212765864 307.4133333333333,5.244680851063895 307.84999999999997,5.861702127659518 308.28666666666663,6.478723404255361 308.7233333333333,6.478723404255361 309.15999999999997,5.861702127659518 309.59666666666664,6.17021276595733 310.0333333333333,5.244680851063895 310.46999999999997,4.627659574468052 310.90666666666664,5.244680851063895 311.3433333333333,5.861702127659518 311.78,6.17021276595733 312.21666666666664,6.478723404255361 312.6533333333333,6.478723404255361 313.09,6.478723404255361 313.52666666666664,6.478723404255361 313.9633333333333,6.17021276595733 314.4,12.957446808510722 314.83666666666664,13.882978723404158 315.2733333333333,13.882978723404158 315.71,11.723404255319036 316.14666666666665,10.180851063829758 316.5833333333333,10.180851063829758 317.02,10.180851063829758 317.45666666666665,11.723404255319036 317.8933333333333,8.63829787234048 318.33,8.63829787234048 318.76666666666665,8.946808510638293 319.2033333333333,8.63829787234048 319.64,8.946808510638293 320.07666666666665,9.255319148936104 320.5133333333333,8.021276595744638 320.95,8.946808510638293 321.38666666666666,9.563829787233916 321.8233333333333,9.563829787233916 322.26,9.872340425531947 322.69666666666666,9.563829787233916 323.1333333333333,9.255319148936104 323.57,12.957446808510722 324.00666666666666,12.648936170212691 324.4433333333333,12.031914893617067 324.88,10.7978723404256 325.31666666666666,11.723404255319036 325.75333333333333,12.34042553191488 326.19,12.031914893617067 326.62666666666667,12.031914893617067 327.06333333333333,11.723404255319036 327.5,11.106382978723413 327.93666666666667,9.872340425531947 328.37333333333333,9.563829787233916 328.81,8.946808510638293 329.24666666666667,8.946808510638293 329.68333333333334,7.712765957446827 330.12,8.021276595744638 330.5566666666667,8.021276595744638 330.99333333333334,7.712765957446827 331.43,7.404255319148796 331.8666666666667,6.787234042553172 332.30333333333334,6.787234042553172 332.74,6.478723404255361 333.1766666666667,6.17021276595733 333.61333333333334,7.712765957446827 334.05,7.404255319148796 334.4866666666667,7.712765957446827 334.92333333333335,8.946808510638293 335.36,9.255319148936104 335.7966666666666,8.63829787234048 336.2333333333333,8.63829787234048 336.66999999999996,9.872340425531947 337.1066666666666,9.563829787233916 337.5433333333333,9.255319148936104 337.97999999999996,10.489361702127571 338.41666666666663,10.7978723404256 338.8533333333333,12.648936170212691 339.28999999999996,11.723404255319036 339.72666666666663,12.957446808510722 340.1633333333333,12.34042553191488 340.59999999999997,12.031914893617067 341.03666666666663,10.489361702127571 341.4733333333333,11.106382978723413 341.90999999999997,10.489361702127571 342.34666666666664,11.414893617021225 342.7833333333333,11.723404255319036 343.21999999999997,11.723404255319036 343.65666666666664,11.414893617021225 344.0933333333333,11.106382978723413 344.53,10.180851063829758 344.96666666666664,9.563829787233916 345.4033333333333,9.255319148936104 345.84,6.478723404255361 346.27666666666664,8.021276595744638 346.7133333333333,8.63829787234048 347.15,8.63829787234048 347.58666666666664,8.021276595744638 348.0233333333333,8.32978723404245 348.46,8.946808510638293 348.89666666666665,8.946808510638293 349.3333333333333,9.872340425531947 349.77,15.425531914893654 350.20666666666665,12.957446808510722 350.6433333333333,16.042553191489276 351.08,12.34042553191488 351.51666666666665,14.191489361702187 351.9533333333333,12.34042553191488 352.39,12.957446808510722 352.82666666666665,13.265957446808534 353.2633333333333,11.414893617021225 353.7,12.031914893617067 354.13666666666666,13.882978723404158 354.5733333333333,13.882978723404158 355.01,13.882978723404158 355.44666666666666,13.882978723404158 355.8833333333333,13.882978723404158 356.32,13.574468085106345 356.75666666666666,12.957446808510722 357.1933333333333,13.265957446808534 357.63,12.648936170212691 358.06666666666666,12.648936170212691 358.50333333333333,12.031914893617067 358.94,12.34042553191488 359.37666666666667,12.031914893617067 359.81333333333333,12.648936170212691 360.25,12.648936170212691 360.68666666666667,12.957446808510722 361.12333333333333,12.648936170212691 361.56,12.34042553191488 361.99666666666667,11.723404255319036 362.43333333333334,11.723404255319036 362.87,9.872340425531947 363.3066666666667,10.489361702127571 363.74333333333334,10.180851063829758 364.18,10.7978723404256 364.6166666666667,10.7978723404256 365.05333333333334,10.7978723404256 365.49,10.489361702127571 365.9266666666667,10.7978723404256 366.36333333333334,10.180851063829758 366.8,10.180851063829758 367.2366666666667,10.489361702127571 367.6733333333333,11.106382978723413 368.10999999999996,12.031914893617067 368.5466666666666,10.7978723404256 368.9833333333333,10.7978723404256 369.41999999999996,10.489361702127571 369.8566666666666,11.723404255319036 370.2933333333333,11.723404255319036 370.72999999999996,12.34042553191488 371.16666666666663,12.648936170212691 371.6033333333333,12.34042553191488 372.03999999999996,12.957446808510722 372.47666666666663,13.882978723404158 372.9133333333333,13.574468085106345 373.34999999999997,13.882978723404158 373.78666666666663,13.574468085106345 374.2233333333333,14.191489361702187 374.65999999999997,14.191489361702187 375.09666666666664,14.5 375.5333333333333,16.351063829787307 375.96999999999997,16.351063829787307 376.40666666666664,17.276595744680744 376.8433333333333,17.276595744680744 377.28,18.81914893617024 377.71666666666664,18.51063829787243 378.1533333333333,18.81914893617024 378.59,19.12765957446805 379.02666666666664,19.12765957446805 379.4633333333333,18.51063829787243 379.9,18.81914893617024 380.33666666666664,19.12765957446805 380.7733333333333,18.51063829787243 381.21,19.12765957446805 381.64666666666665,19.744680851063894 382.0833333333333,20.97872340425536 382.52,20.053191489361705 382.95666666666665,29 383.3933333333333,25.914893617021225 383.83,23.13829787234048 384.26666666666665,22.521276595744638 384.7033333333333,23.44680851063829 385.14,21.904255319149016 385.57666666666665,22.212765957446827 386.0133333333333,21.904255319149016 386.45,22.82978723404245 386.88666666666666,22.82978723404245 387.3233333333333,23.44680851063829 387.76,23.44680851063829 388.19666666666666,23.755319148936103 388.6333333333333,24.68085106382976 389.07,24.372340425531945 389.50666666666666,23.755319148936103 389.9433333333333,24.063829787234134 390.38,24.063829787234134 390.81666666666666,24.68085106382976 391.25333333333333,24.68085106382976 391.69,24.063829787234134 392.12666666666667,24.063829787234134 392.56333333333333,24.063829787234134 393,23.755319148936103 393.43666666666667,22.82978723404245 393.87333333333333,22.212765957446827 394.31,22.212765957446827 394.74666666666667,22.212765957446827 395.18333333333334,22.82978723404245 395.62,22.82978723404245 396.0566666666667,21.287234042553173 396.49333333333334,21.595744680850984 396.93,21.595744680850984 397.3666666666667,21.287234042553173 397.80333333333334,21.904255319149016 398.24,21.287234042553173 398.6766666666667,22.521276595744638 399.11333333333334,22.521276595744638 399.54999999999995,21.595744680850984 399.9866666666666,21.904255319149016 400.4233333333333,20.97872340425536 400.85999999999996,21.595744680850984 401.2966666666666,21.287234042553173 401.7333333333333,20.670212765957547 402.16999999999996,21.904255319149016 402.6066666666666,21.287234042553173 403.0433333333333,21.904255319149016 403.47999999999996,21.595744680850984 403.91666666666663,22.521276595744638 404.3533333333333,21.904255319149016 404.78999999999996,21.904255319149016 405.22666666666663,21.287234042553173 405.6633333333333,21.287234042553173 406.09999999999997,20.670212765957547 406.53666666666663,21.595744680850984 406.9733333333333,20.670212765957547 407.40999999999997,20.97872340425536 407.84666666666664,20.361702127659516 408.2833333333333,20.361702127659516 408.71999999999997,20.670212765957547 409.15666666666664,19.744680851063894 409.5933333333333,20.053191489361705 410.03,19.12765957446805 410.46666666666664,19.744680851063894 410.9033333333333,20.053191489361705 411.34,20.053191489361705 411.77666666666664,20.053191489361705 412.2133333333333,19.436170212765862 412.65,18.51063829787243 413.08666666666664,19.12765957446805 413.5233333333333,19.12765957446805 413.96,19.12765957446805 414.39666666666665,18.51063829787243 414.8333333333333,19.436170212765862 415.27,19.12765957446805 415.70666666666665,19.436170212765862 416.1433333333333,19.436170212765862 416.58,20.361702127659516 417.01666666666665,20.670212765957547 417.4533333333333,20.361702127659516 417.89,20.97872340425536 418.32666666666665,20.97872340425536 418.7633333333333,22.82978723404245 419.2,24.063829787234134 419.63666666666666,20.97872340425536 420.0733333333333,22.82978723404245 420.51,25.606382978723413 420.94666666666666,22.82978723404245 421.3833333333333,21.595744680850984 421.82,20.97872340425536 422.25666666666666,20.670212765957547 422.6933333333333,20.361702127659516 423.13,20.97872340425536 423.56666666666666,21.595744680850984 424.00333333333333,20.053191489361705 424.44,20.361702127659516 424.87666666666667,20.670212765957547 425.31333333333333,20.670212765957547 425.75,19.744680851063894 426.18666666666667,21.595744680850984 426.62333333333333,20.97872340425536 427.06,19.744680851063894 427.49666666666667,20.361702127659516 427.93333333333334,20.97872340425536 428.37,23.44680851063829 428.8066666666667,22.82978723404245 429.24333333333334,22.82978723404245 429.68,22.521276595744638 430.1166666666667,17.893617021276587 430.55333333333334,18.51063829787243 430.99,18.51063829787243 431.4266666666666,18.81914893617024 431.8633333333333,17.893617021276587 432.29999999999995,16.351063829787307 432.7366666666666,16.65957446808512 433.1733333333333,14.808510638297811 433.60999999999996,16.65957446808512 434.0466666666666,16.042553191489276 434.4833333333333,15.425531914893654 434.91999999999996,13.574468085106345 435.3566666666666,14.5 435.7933333333333,14.191489361702187 436.22999999999996,14.191489361702187 436.66666666666663,12.34042553191488 437.1033333333333,12.031914893617067 437.53999999999996,11.414893617021225 437.97666666666663,9.563829787233916 438.4133333333333,9.255319148936104 438.84999999999997,8.946808510638293 439.28666666666663,10.7978723404256 439.7233333333333,10.489361702127571 440.15999999999997,10.7978723404256 440.59666666666664,10.7978723404256 441.0333333333333,11.106382978723413 441.46999999999997,10.7978723404256 441.90666666666664,10.180851063829758 442.3433333333333,10.180851063829758 442.78,9.563829787233916 443.21666666666664,10.7978723404256 443.6533333333333,9.563829787233916 444.09,8.63829787234048 444.52666666666664,8.946808510638293 444.9633333333333,9.255319148936104 445.4,8.021276595744638 445.83666666666664,7.404255319148796 446.2733333333333,7.712765957446827 446.71,6.478723404255361 447.14666666666665,8.32978723404245 447.5833333333333,7.404255319148796 448.02,8.63829787234048 448.45666666666665,9.255319148936104 448.8933333333333,8.32978723404245 449.33,6.787234042553172 449.76666666666665,7.712765957446827 450.2033333333333,9.255319148936104 450.64,8.32978723404245 451.07666666666665,11.414893617021225 451.5133333333333,10.180851063829758 451.95,12.031914893617067 452.38666666666666,10.489361702127571 452.8233333333333,10.489361702127571 453.26,11.723404255319036 453.69666666666666,12.031914893617067 454.1333333333333,11.723404255319036 454.57,12.031914893617067 455.00666666666666,11.106382978723413 455.4433333333333,12.648936170212691 455.88,11.414893617021225 456.31666666666666,12.34042553191488 456.75333333333333,13.265957446808534 457.19,14.5 457.62666666666667,14.191489361702187 458.06333333333333,13.574468085106345 458.5,11.723404255319036 458.93666666666667,11.723404255319036 459.37333333333333,12.031914893617067 459.81,10.7978723404256 460.24666666666667,10.180851063829758 460.68333333333334,9.872340425531947 461.12,9.255319148936104 461.5566666666667,9.563829787233916 461.99333333333334,7.712765957446827 462.43,8.63829787234048 462.8666666666667,8.946808510638293 463.3033333333333,9.563829787233916 463.73999999999995,9.872340425531947 464.1766666666666,10.180851063829758 464.6133333333333,8.946808510638293 465.04999999999995,8.63829787234048 465.4866666666666,7.404255319148796 465.9233333333333,8.32978723404245 466.35999999999996,8.021276595744638 466.7966666666666,8.32978723404245 467.2333333333333,7.404255319148796 467.66999999999996,7.404255319148796 468.1066666666666,6.478723404255361 468.5433333333333,6.17021276595733 468.97999999999996,8.32978723404245 469.41666666666663,8.021276595744638 469.8533333333333,8.021276595744638 470.28999999999996,8.021276595744638 470.72666666666663,6.478723404255361 471.1633333333333,7.712765957446827 471.59999999999997,8.021276595744638 472.03666666666663,9.563829787233916 472.4733333333333,8.946808510638293 472.90999999999997,7.404255319148796 473.34666666666664,8.021276595744638 473.7833333333333,7.095744680850984 474.21999999999997,7.404255319148796 474.65666666666664,7.712765957446827 475.0933333333333,8.63829787234048 475.53,8.63829787234048 475.96666666666664,8.32978723404245 476.4033333333333,9.563829787233916 476.84,10.7978723404256 477.27666666666664,11.106382978723413 477.7133333333333,11.723404255319036 478.15,10.7978723404256 478.58666666666664,9.872340425531947 479.0233333333333,10.489361702127571 479.46,11.106382978723413 479.89666666666665,11.414893617021225 480.3333333333333,10.489361702127571 480.77,10.180851063829758 481.20666666666665,10.489361702127571 481.6433333333333,10.7978723404256 482.08,10.489361702127571 482.51666666666665,10.7978723404256 482.9533333333333,11.414893617021225 483.39,11.106382978723413 483.82666666666665,11.106382978723413 484.2633333333333,10.7978723404256 484.7,11.414893617021225 485.13666666666666,11.414893617021225 485.5733333333333,11.414893617021225 486.01,11.723404255319036 486.44666666666666,12.031914893617067 486.8833333333333,12.34042553191488 487.32,12.031914893617067 487.75666666666666,12.957446808510722 488.1933333333333,12.957446808510722 488.63,13.265957446808534 489.06666666666666,12.957446808510722 489.50333333333333,12.648936170212691 489.94,12.957446808510722 490.37666666666667,12.957446808510722 490.81333333333333,13.265957446808534 491.25,12.957446808510722 491.68666666666667,12.648936170212691 492.12333333333333,12.957446808510722 492.56,12.648936170212691 492.99666666666667,12.031914893617067 493.43333333333334,13.574468085106345 493.87,13.574468085106345 494.3066666666667,14.191489361702187 494.74333333333334,13.265957446808534 495.17999999999995,13.574468085106345 495.6166666666666,14.5 496.0533333333333,13.882978723404158 496.48999999999995,13.574468085106345 496.9266666666666,13.574468085106345 497.3633333333333,13.265957446808534 497.79999999999995,12.34042553191488 498.2366666666666,13.882978723404158 498.6733333333333,12.957446808510722 499.10999999999996,12.031914893617067 499.5466666666666,14.191489361702187 499.9833333333333,13.882978723404158 500.41999999999996,13.882978723404158 500.8566666666666,13.882978723404158 501.2933333333333,13.574468085106345 501.72999999999996,13.574468085106345 502.16666666666663,12.957446808510722 502.6033333333333,12.957446808510722 503.03999999999996,13.574468085106345 503.47666666666663,12.648936170212691 503.9133333333333,12.031914893617067 504.34999999999997,11.106382978723413 504.78666666666663,10.7978723404256 505.2233333333333,12.34042553191488 505.65999999999997,12.34042553191488 506.09666666666664,12.031914893617067 506.5333333333333,12.031914893617067 506.96999999999997,10.489361702127571 507.40666666666664,11.414893617021225 507.8433333333333,12.031914893617067 508.28,12.34042553191488 508.71666666666664,12.34042553191488 509.1533333333333,12.34042553191488 509.59,12.031914893617067 510.02666666666664,11.414893617021225 510.4633333333333,11.106382978723413 510.9,12.031914893617067 511.33666666666664,12.031914893617067 511.7733333333333,12.957446808510722 512.2099999999999,13.574468085106345 512.6466666666666,12.34042553191488 513.0833333333333,12.957446808510722 513.52,13.265957446808534 513.9566666666666,12.957446808510722 514.3933333333333,12.957446808510722 514.8299999999999,12.957446808510722 515.2666666666667,13.574468085106345 515.7033333333333,13.574468085106345 516.14,13.574468085106345 516.5766666666666,14.5 517.0133333333333,13.882978723404158 517.4499999999999,13.574468085106345 517.8866666666667,14.5 518.3233333333333,14.191489361702187 518.76,12.34042553191488 519.1966666666666,13.882978723404158 519.6333333333333,13.882978723404158 520.0699999999999,14.5 520.5066666666667,14.5 520.9433333333333,13.882978723404158 521.38,13.882978723404158 521.8166666666666,13.882978723404158 522.2533333333333,12.957446808510722 522.6899999999999,14.808510638297811 523.1266666666667,13.574468085106345 523.5633333333333,13.574468085106345 524,13.574468085106345 524,13.574468085106345 524,47 0,47 0,22.82978723404245" style="fill:url(#market_chart_bar);stroke: none;" stroke-dasharray="" stroke-width="undefined"></polyline>
            
        <polyline points="0,22.82978723404245 0.43666666666666665,25.914893617021225 0.8733333333333333,23.755319148936103 1.31,22.521276595744638 1.7466666666666666,22.212765957446827 2.183333333333333,21.904255319149016 2.62,21.595744680850984 3.0566666666666666,20.670212765957547 3.493333333333333,19.436170212765862 3.9299999999999997,19.436170212765862 4.366666666666666,20.670212765957547 4.803333333333333,21.287234042553173 5.24,20.97872340425536 5.676666666666667,21.904255319149016 6.113333333333333,22.521276595744638 6.55,22.521276595744638 6.986666666666666,21.287234042553173 7.423333333333333,20.053191489361705 7.859999999999999,19.12765957446805 8.296666666666667,20.97872340425536 8.733333333333333,20.97872340425536 9.17,20.670212765957547 9.606666666666666,21.904255319149016 10.043333333333333,21.595744680850984 10.48,19.436170212765862 10.916666666666666,20.053191489361705 11.353333333333333,19.12765957446805 11.79,16.968085106382933 12.226666666666667,15.734042553191465 12.663333333333332,15.425531914893654 13.1,16.351063829787307 13.536666666666665,15.734042553191465 13.973333333333333,16.042553191489276 14.41,15.734042553191465 14.846666666666666,15.734042553191465 15.283333333333333,17.276595744680744 15.719999999999999,17.893617021276587 16.156666666666666,16.968085106382933 16.593333333333334,16.968085106382933 17.029999999999998,17.276595744680744 17.466666666666665,16.351063829787307 17.903333333333332,16.65957446808512 18.34,15.425531914893654 18.776666666666667,15.734042553191465 19.21333333333333,15.117021276595622 19.65,14.808510638297811 20.086666666666666,14.808510638297811 20.523333333333333,15.734042553191465 20.96,16.351063829787307 21.396666666666665,17.585106382978775 21.833333333333332,16.65957446808512 22.27,15.425531914893654 22.706666666666667,17.276595744680744 23.14333333333333,16.968085106382933 23.58,16.65957446808512 24.016666666666666,16.65957446808512 24.453333333333333,15.734042553191465 24.89,16.351063829787307 25.326666666666664,16.042553191489276 25.763333333333332,15.425531914893654 26.2,15.425531914893654 26.636666666666667,13.574468085106345 27.07333333333333,14.191489361702187 27.509999999999998,16.65957446808512 27.946666666666665,16.351063829787307 28.383333333333333,17.276595744680744 28.82,16.351063829787307 29.256666666666664,16.042553191489276 29.69333333333333,17.276595744680744 30.13,17.276595744680744 30.566666666666666,16.351063829787307 31.00333333333333,16.65957446808512 31.439999999999998,16.351063829787307 31.876666666666665,15.425531914893654 32.31333333333333,16.042553191489276 32.75,17.585106382978775 33.18666666666667,18.202127659574398 33.623333333333335,16.968085106382933 34.059999999999995,16.65957446808512 34.49666666666666,16.65957446808512 34.93333333333333,17.276595744680744 35.37,17.893617021276587 35.806666666666665,18.81914893617024 36.24333333333333,17.276595744680744 36.68,18.202127659574398 37.11666666666667,19.436170212765862 37.553333333333335,19.744680851063894 37.989999999999995,19.436170212765862 38.42666666666666,20.670212765957547 38.86333333333333,20.053191489361705 39.3,19.744680851063894 39.736666666666665,21.904255319149016 40.17333333333333,18.51063829787243 40.61,19.12765957446805 41.04666666666667,20.053191489361705 41.483333333333334,21.287234042553173 41.92,22.82978723404245 42.35666666666666,23.13829787234048 42.79333333333333,24.063829787234134 43.23,25.297872340425602 43.666666666666664,23.13829787234048 44.10333333333333,23.44680851063829 44.54,24.063829787234134 44.97666666666667,24.372340425531945 45.413333333333334,23.755319148936103 45.85,24.372340425531945 46.28666666666666,25.606382978723413 46.72333333333333,27.14893617021269 47.16,24.68085106382976 47.596666666666664,24.063829787234134 48.03333333333333,23.755319148936103 48.47,23.44680851063829 48.906666666666666,23.13829787234048 49.343333333333334,23.13829787234048 49.78,22.82978723404245 50.21666666666666,22.212765957446827 50.65333333333333,21.595744680850984 51.089999999999996,22.521276595744638 51.526666666666664,22.521276595744638 51.96333333333333,22.521276595744638 52.4,23.44680851063829 52.836666666666666,21.595744680850984 53.27333333333333,23.13829787234048 53.71,23.755319148936103 54.14666666666666,24.98936170212757 54.58333333333333,24.372340425531945 55.019999999999996,23.755319148936103 55.45666666666666,23.44680851063829 55.89333333333333,23.44680851063829 56.33,22.82978723404245 56.766666666666666,22.82978723404245 57.20333333333333,23.13829787234048 57.64,23.44680851063829 58.07666666666666,23.755319148936103 58.51333333333333,22.521276595744638 58.949999999999996,21.287234042553173 59.38666666666666,21.595744680850984 59.82333333333333,23.755319148936103 60.26,23.44680851063829 60.696666666666665,22.521276595744638 61.13333333333333,22.212765957446827 61.57,22.82978723404245 62.00666666666666,21.904255319149016 62.44333333333333,21.595744680850984 62.879999999999995,21.904255319149016 63.31666666666666,21.595744680850984 63.75333333333333,21.595744680850984 64.19,20.97872340425536 64.62666666666667,20.97872340425536 65.06333333333333,20.97872340425536 65.5,20.053191489361705 65.93666666666667,20.361702127659516 66.37333333333333,19.744680851063894 66.81,19.436170212765862 67.24666666666667,19.436170212765862 67.68333333333334,19.12765957446805 68.11999999999999,19.12765957446805 68.55666666666666,19.744680851063894 68.99333333333333,17.276595744680744 69.42999999999999,16.042553191489276 69.86666666666666,15.425531914893654 70.30333333333333,15.734042553191465 70.74,16.351063829787307 71.17666666666666,15.734042553191465 71.61333333333333,15.425531914893654 72.05,14.808510638297811 72.48666666666666,14.5 72.92333333333333,13.265957446808534 73.36,12.648936170212691 73.79666666666667,14.191489361702187 74.23333333333333,16.042553191489276 74.67,15.425531914893654 75.10666666666667,15.425531914893654 75.54333333333334,14.808510638297811 75.97999999999999,14.5 76.41666666666666,14.808510638297811 76.85333333333332,15.117021276595622 77.28999999999999,15.425531914893654 77.72666666666666,15.117021276595622 78.16333333333333,15.734042553191465 78.6,15.734042553191465 79.03666666666666,15.425531914893654 79.47333333333333,15.425531914893654 79.91,15.425531914893654 80.34666666666666,14.191489361702187 80.78333333333333,15.117021276595622 81.22,12.648936170212691 81.65666666666667,12.031914893617067 82.09333333333333,10.489361702127571 82.53,11.106382978723413 82.96666666666667,12.34042553191488 83.40333333333334,10.7978723404256 83.84,10.489361702127571 84.27666666666666,9.872340425531947 84.71333333333332,10.180851063829758 85.14999999999999,11.106382978723413 85.58666666666666,11.106382978723413 86.02333333333333,11.414893617021225 86.46,10.7978723404256 86.89666666666666,10.180851063829758 87.33333333333333,11.723404255319036 87.77,13.265957446808534 88.20666666666666,12.648936170212691 88.64333333333333,13.265957446808534 89.08,13.574468085106345 89.51666666666667,14.5 89.95333333333333,14.808510638297811 90.39,15.425531914893654 90.82666666666667,13.574468085106345 91.26333333333334,14.5 91.7,14.191489361702187 92.13666666666666,14.808510638297811 92.57333333333332,14.808510638297811 93.00999999999999,14.191489361702187 93.44666666666666,14.191489361702187 93.88333333333333,12.957446808510722 94.32,13.882978723404158 94.75666666666666,12.34042553191488 95.19333333333333,12.648936170212691 95.63,12.957446808510722 96.06666666666666,12.957446808510722 96.50333333333333,12.34042553191488 96.94,11.106382978723413 97.37666666666667,10.489361702127571 97.81333333333333,10.7978723404256 98.25,11.723404255319036 98.68666666666667,13.265957446808534 99.12333333333333,12.031914893617067 99.56,11.723404255319036 99.99666666666666,12.648936170212691 100.43333333333332,12.648936170212691 100.86999999999999,12.648936170212691 101.30666666666666,12.031914893617067 101.74333333333333,12.031914893617067 102.17999999999999,11.106382978723413 102.61666666666666,10.7978723404256 103.05333333333333,11.106382978723413 103.49,11.106382978723413 103.92666666666666,11.106382978723413 104.36333333333333,10.180851063829758 104.8,9.872340425531947 105.23666666666666,11.414893617021225 105.67333333333333,11.106382978723413 106.11,11.106382978723413 106.54666666666667,8.63829787234048 106.98333333333333,7.095744680850984 107.42,8.021276595744638 107.85666666666665,11.414893617021225 108.29333333333332,13.574468085106345 108.72999999999999,12.957446808510722 109.16666666666666,14.5 109.60333333333332,13.882978723404158 110.03999999999999,14.5 110.47666666666666,12.957446808510722 110.91333333333333,13.265957446808534 111.35,11.414893617021225 111.78666666666666,13.265957446808534 112.22333333333333,12.957446808510722 112.66,12.031914893617067 113.09666666666666,11.106382978723413 113.53333333333333,12.957446808510722 113.97,10.7978723404256 114.40666666666667,9.872340425531947 114.84333333333333,6.787234042553172 115.28,5.244680851063895 115.71666666666667,5.553191489361707 116.15333333333332,4.936170212765864 116.58999999999999,2.468085106382932 117.02666666666666,4.01063829787221 117.46333333333332,6.478723404255361 117.89999999999999,5.244680851063895 118.33666666666666,4.31914893617024 118.77333333333333,5.244680851063895 119.21,3.702127659574398 119.64666666666666,2.15957446808512 120.08333333333333,4.01063829787221 120.52,3.393617021276586 120.95666666666666,0 121.39333333333333,0 121.83,1.5425531914892776 122.26666666666667,2.15957446808512 122.70333333333333,2.15957446808512 123.14,2.15957446808512 123.57666666666667,2.7765957446807437 124.01333333333332,2.15957446808512 124.44999999999999,4.31914893617024 124.88666666666666,3.0851063829787746 125.32333333333332,3.702127659574398 125.75999999999999,4.627659574468052 126.19666666666666,4.627659574468052 126.63333333333333,5.861702127659518 127.07,6.787234042553172 127.50666666666666,8.021276595744638 127.94333333333333,8.32978723404245 128.38,8.946808510638293 128.81666666666666,10.180851063829758 129.25333333333333,10.7978723404256 129.69,9.255319148936104 130.12666666666667,10.180851063829758 130.56333333333333,8.021276595744638 131,8.946808510638293 131.43666666666667,9.255319148936104 131.87333333333333,9.563829787233916 132.31,10.180851063829758 132.74666666666667,10.180851063829758 133.18333333333334,10.7978723404256 133.62,10.489361702127571 134.05666666666667,11.106382978723413 134.49333333333334,11.723404255319036 134.93,9.563829787233916 135.36666666666667,9.563829787233916 135.8033333333333,10.180851063829758 136.23999999999998,10.489361702127571 136.67666666666665,10.180851063829758 137.11333333333332,9.872340425531947 137.54999999999998,9.255319148936104 137.98666666666665,9.872340425531947 138.42333333333332,8.63829787234048 138.85999999999999,8.63829787234048 139.29666666666665,9.563829787233916 139.73333333333332,9.563829787233916 140.17,9.255319148936104 140.60666666666665,9.872340425531947 141.04333333333332,10.489361702127571 141.48,10.7978723404256 141.91666666666666,10.7978723404256 142.35333333333332,10.489361702127571 142.79,11.106382978723413 143.22666666666666,10.7978723404256 143.66333333333333,11.723404255319036 144.1,12.957446808510722 144.53666666666666,12.34042553191488 144.97333333333333,12.031914893617067 145.41,12.031914893617067 145.84666666666666,12.031914893617067 146.28333333333333,10.489361702127571 146.72,11.723404255319036 147.15666666666667,12.031914893617067 147.59333333333333,11.723404255319036 148.03,12.34042553191488 148.46666666666667,11.106382978723413 148.90333333333334,9.872340425531947 149.34,9.563829787233916 149.77666666666667,9.563829787233916 150.21333333333334,9.563829787233916 150.65,10.180851063829758 151.08666666666667,12.34042553191488 151.52333333333334,12.957446808510722 151.95999999999998,14.191489361702187 152.39666666666665,13.574468085106345 152.83333333333331,12.031914893617067 153.26999999999998,12.648936170212691 153.70666666666665,11.723404255319036 154.14333333333332,12.957446808510722 154.57999999999998,12.957446808510722 155.01666666666665,13.882978723404158 155.45333333333332,13.882978723404158 155.89,13.574468085106345 156.32666666666665,13.574468085106345 156.76333333333332,13.882978723404158 157.2,16.968085106382933 157.63666666666666,16.65957446808512 158.07333333333332,15.734042553191465 158.51,15.734042553191465 158.94666666666666,16.042553191489276 159.38333333333333,16.042553191489276 159.82,15.425531914893654 160.25666666666666,15.425531914893654 160.69333333333333,16.042553191489276 161.13,16.042553191489276 161.56666666666666,16.042553191489276 162.00333333333333,14.5 162.44,15.117021276595622 162.87666666666667,13.882978723404158 163.31333333333333,14.808510638297811 163.75,13.574468085106345 164.18666666666667,13.882978723404158 164.62333333333333,13.882978723404158 165.06,13.882978723404158 165.49666666666667,14.191489361702187 165.93333333333334,12.648936170212691 166.37,10.489361702127571 166.80666666666667,11.723404255319036 167.24333333333334,11.414893617021225 167.68,10.489361702127571 168.11666666666665,10.180851063829758 168.5533333333333,10.7978723404256 168.98999999999998,10.180851063829758 169.42666666666665,10.489361702127571 169.86333333333332,9.872340425531947 170.29999999999998,12.648936170212691 170.73666666666665,12.34042553191488 171.17333333333332,12.648936170212691 171.60999999999999,10.489361702127571 172.04666666666665,8.63829787234048 172.48333333333332,9.563829787233916 172.92,8.021276595744638 173.35666666666665,8.021276595744638 173.79333333333332,8.021276595744638 174.23,5.861702127659518 174.66666666666666,5.861702127659518 175.10333333333332,6.787234042553172 175.54,4.936170212765864 175.97666666666666,8.021276595744638 176.41333333333333,8.63829787234048 176.85,8.32978723404245 177.28666666666666,8.021276595744638 177.72333333333333,8.946808510638293 178.16,10.489361702127571 178.59666666666666,8.946808510638293 179.03333333333333,8.63829787234048 179.47,8.021276595744638 179.90666666666667,8.946808510638293 180.34333333333333,8.63829787234048 180.78,9.563829787233916 181.21666666666667,9.255319148936104 181.65333333333334,10.489361702127571 182.09,11.106382978723413 182.52666666666667,10.7978723404256 182.96333333333334,11.106382978723413 183.4,11.106382978723413 183.83666666666664,10.489361702127571 184.2733333333333,11.414893617021225 184.70999999999998,11.106382978723413 185.14666666666665,11.414893617021225 185.58333333333331,11.414893617021225 186.01999999999998,11.723404255319036 186.45666666666665,12.031914893617067 186.89333333333332,13.882978723404158 187.32999999999998,14.191489361702187 187.76666666666665,14.5 188.20333333333332,14.808510638297811 188.64,14.191489361702187 189.07666666666665,14.5 189.51333333333332,14.191489361702187 189.95,14.5 190.38666666666666,14.5 190.82333333333332,14.808510638297811 191.26,13.882978723404158 191.69666666666666,14.808510638297811 192.13333333333333,14.808510638297811 192.57,14.191489361702187 193.00666666666666,14.808510638297811 193.44333333333333,14.5 193.88,14.808510638297811 194.31666666666666,14.808510638297811 194.75333333333333,15.117021276595622 195.19,12.031914893617067 195.62666666666667,12.031914893617067 196.06333333333333,11.106382978723413 196.5,11.723404255319036 196.93666666666667,11.414893617021225 197.37333333333333,12.34042553191488 197.81,13.265957446808534 198.24666666666667,11.414893617021225 198.68333333333334,10.7978723404256 199.12,11.723404255319036 199.55666666666667,13.574468085106345 199.9933333333333,12.031914893617067 200.42999999999998,12.957446808510722 200.86666666666665,12.031914893617067 201.3033333333333,12.957446808510722 201.73999999999998,15.117021276595622 202.17666666666665,16.65957446808512 202.61333333333332,17.585106382978775 203.04999999999998,16.042553191489276 203.48666666666665,17.276595744680744 203.92333333333332,16.351063829787307 204.35999999999999,17.276595744680744 204.79666666666665,18.202127659574398 205.23333333333332,18.51063829787243 205.67,17.893617021276587 206.10666666666665,16.968085106382933 206.54333333333332,17.585106382978775 206.98,17.585106382978775 207.41666666666666,19.436170212765862 207.85333333333332,19.436170212765862 208.29,19.436170212765862 208.72666666666666,19.436170212765862 209.16333333333333,17.893617021276587 209.6,13.574468085106345 210.03666666666666,14.5 210.47333333333333,14.5 210.91,14.5 211.34666666666666,17.893617021276587 211.78333333333333,19.12765957446805 212.22,18.81914893617024 212.65666666666667,19.436170212765862 213.09333333333333,17.893617021276587 213.53,17.276595744680744 213.96666666666667,16.968085106382933 214.40333333333334,18.51063829787243 214.84,18.81914893617024 215.27666666666667,18.202127659574398 215.7133333333333,17.276595744680744 216.14999999999998,16.351063829787307 216.58666666666664,17.893617021276587 217.0233333333333,16.351063829787307 217.45999999999998,15.425531914893654 217.89666666666665,16.351063829787307 218.33333333333331,16.968085106382933 218.76999999999998,15.117021276595622 219.20666666666665,15.117021276595622 219.64333333333332,15.117021276595622 220.07999999999998,16.65957446808512 220.51666666666665,17.893617021276587 220.95333333333332,18.202127659574398 221.39,15.734042553191465 221.82666666666665,17.585106382978775 222.26333333333332,16.351063829787307 222.7,16.968085106382933 223.13666666666666,16.968085106382933 223.57333333333332,19.744680851063894 224.01,20.053191489361705 224.44666666666666,20.053191489361705 224.88333333333333,17.893617021276587 225.32,17.276595744680744 225.75666666666666,17.276595744680744 226.19333333333333,17.276595744680744 226.63,18.202127659574398 227.06666666666666,18.51063829787243 227.50333333333333,19.744680851063894 227.94,19.744680851063894 228.37666666666667,19.744680851063894 228.81333333333333,19.436170212765862 229.25,19.744680851063894 229.68666666666667,19.436170212765862 230.12333333333333,19.12765957446805 230.56,18.81914893617024 230.99666666666667,18.202127659574398 231.43333333333334,17.585106382978775 231.86999999999998,16.968085106382933 232.30666666666664,18.51063829787243 232.7433333333333,18.51063829787243 233.17999999999998,18.202127659574398 233.61666666666665,17.276595744680744 234.0533333333333,16.65957446808512 234.48999999999998,16.042553191489276 234.92666666666665,14.5 235.36333333333332,14.5 235.79999999999998,14.808510638297811 236.23666666666665,15.117021276595622 236.67333333333332,15.425531914893654 237.10999999999999,16.042553191489276 237.54666666666665,16.042553191489276 237.98333333333332,13.882978723404158 238.42,14.191489361702187 238.85666666666665,13.574468085106345 239.29333333333332,13.574468085106345 239.73,11.106382978723413 240.16666666666666,10.489361702127571 240.60333333333332,10.489361702127571 241.04,10.180851063829758 241.47666666666666,8.63829787234048 241.91333333333333,8.946808510638293 242.35,10.180851063829758 242.78666666666666,10.180851063829758 243.22333333333333,11.106382978723413 243.66,10.7978723404256 244.09666666666666,10.180851063829758 244.53333333333333,10.489361702127571 244.97,8.021276595744638 245.40666666666667,5.553191489361707 245.84333333333333,6.787234042553172 246.28,6.17021276595733 246.71666666666667,5.861702127659518 247.15333333333334,6.787234042553172 247.58999999999997,7.095744680850984 248.02666666666664,7.712765957446827 248.4633333333333,8.021276595744638 248.89999999999998,7.095744680850984 249.33666666666664,6.787234042553172 249.7733333333333,7.095744680850984 250.20999999999998,7.095744680850984 250.64666666666665,6.478723404255361 251.08333333333331,6.478723404255361 251.51999999999998,5.553191489361707 251.95666666666665,5.861702127659518 252.39333333333332,6.17021276595733 252.82999999999998,6.478723404255361 253.26666666666665,6.17021276595733 253.70333333333332,5.861702127659518 254.14,5.861702127659518 254.57666666666665,4.01063829787221 255.01333333333332,4.936170212765864 255.45,5.244680851063895 255.88666666666666,4.936170212765864 256.3233333333333,4.936170212765864 256.76,5.553191489361707 257.19666666666666,3.0851063829787746 257.6333333333333,3.702127659574398 258.07,3.393617021276586 258.50666666666666,4.01063829787221 258.9433333333333,4.627659574468052 259.38,4.627659574468052 259.81666666666666,4.31914893617024 260.25333333333333,4.936170212765864 260.69,4.936170212765864 261.12666666666667,5.244680851063895 261.56333333333333,2.7765957446807437 262,0.9255319148936543 262.43666666666667,2.468085106382932 262.87333333333333,1.234042553191466 263.31,0.30851063829781167 263.74666666666667,1.5425531914892776 264.18333333333334,1.5425531914892776 264.62,1.8510638297870894 265.0566666666667,2.15957446808512 265.49333333333334,2.7765957446807437 265.93,3.702127659574398 266.3666666666667,3.393617021276586 266.80333333333334,4.01063829787221 267.24,4.01063829787221 267.6766666666667,4.936170212765864 268.11333333333334,3.0851063829787746 268.55,3.702127659574398 268.9866666666667,4.01063829787221 269.42333333333335,3.702127659574398 269.86,3.393617021276586 270.2966666666667,4.01063829787221 270.73333333333335,4.936170212765864 271.17,3.702127659574398 271.6066666666666,3.0851063829787746 272.0433333333333,4.01063829787221 272.47999999999996,3.702127659574398 272.91666666666663,4.01063829787221 273.3533333333333,4.01063829787221 273.78999999999996,4.01063829787221 274.22666666666663,4.31914893617024 274.6633333333333,4.31914893617024 275.09999999999997,5.244680851063895 275.53666666666663,5.244680851063895 275.9733333333333,5.861702127659518 276.40999999999997,5.861702127659518 276.84666666666664,5.861702127659518 277.2833333333333,6.478723404255361 277.71999999999997,6.478723404255361 278.15666666666664,6.17021276595733 278.5933333333333,6.17021276595733 279.03,6.478723404255361 279.46666666666664,6.17021276595733 279.9033333333333,6.787234042553172 280.34,6.478723404255361 280.77666666666664,7.095744680850984 281.2133333333333,7.095744680850984 281.65,7.404255319148796 282.08666666666664,8.021276595744638 282.5233333333333,8.021276595744638 282.96,7.712765957446827 283.39666666666665,7.404255319148796 283.8333333333333,7.095744680850984 284.27,7.095744680850984 284.70666666666665,5.553191489361707 285.1433333333333,6.17021276595733 285.58,5.861702127659518 286.01666666666665,6.17021276595733 286.4533333333333,6.17021276595733 286.89,6.478723404255361 287.32666666666665,6.787234042553172 287.7633333333333,6.787234042553172 288.2,6.787234042553172 288.63666666666666,6.787234042553172 289.0733333333333,6.478723404255361 289.51,6.478723404255361 289.94666666666666,6.478723404255361 290.3833333333333,7.095744680850984 290.82,7.095744680850984 291.25666666666666,6.478723404255361 291.6933333333333,6.478723404255361 292.13,6.787234042553172 292.56666666666666,6.787234042553172 293.00333333333333,7.095744680850984 293.44,7.404255319148796 293.87666666666667,7.712765957446827 294.31333333333333,8.021276595744638 294.75,8.021276595744638 295.18666666666667,7.712765957446827 295.62333333333333,7.712765957446827 296.06,7.404255319148796 296.49666666666667,7.404255319148796 296.93333333333334,8.32978723404245 297.37,8.32978723404245 297.8066666666667,7.095744680850984 298.24333333333334,7.095744680850984 298.68,7.095744680850984 299.1166666666667,6.478723404255361 299.55333333333334,6.787234042553172 299.99,6.787234042553172 300.4266666666667,6.787234042553172 300.86333333333334,6.787234042553172 301.3,6.478723404255361 301.7366666666667,6.787234042553172 302.17333333333335,6.787234042553172 302.61,6.478723404255361 303.0466666666667,6.478723404255361 303.4833333333333,6.17021276595733 303.91999999999996,3.702127659574398 304.3566666666666,4.01063829787221 304.7933333333333,4.31914893617024 305.22999999999996,4.31914893617024 305.66666666666663,3.393617021276586 306.1033333333333,2.468085106382932 306.53999999999996,4.936170212765864 306.97666666666663,4.936170212765864 307.4133333333333,5.244680851063895 307.84999999999997,5.861702127659518 308.28666666666663,6.478723404255361 308.7233333333333,6.478723404255361 309.15999999999997,5.861702127659518 309.59666666666664,6.17021276595733 310.0333333333333,5.244680851063895 310.46999999999997,4.627659574468052 310.90666666666664,5.244680851063895 311.3433333333333,5.861702127659518 311.78,6.17021276595733 312.21666666666664,6.478723404255361 312.6533333333333,6.478723404255361 313.09,6.478723404255361 313.52666666666664,6.478723404255361 313.9633333333333,6.17021276595733 314.4,12.957446808510722 314.83666666666664,13.882978723404158 315.2733333333333,13.882978723404158 315.71,11.723404255319036 316.14666666666665,10.180851063829758 316.5833333333333,10.180851063829758 317.02,10.180851063829758 317.45666666666665,11.723404255319036 317.8933333333333,8.63829787234048 318.33,8.63829787234048 318.76666666666665,8.946808510638293 319.2033333333333,8.63829787234048 319.64,8.946808510638293 320.07666666666665,9.255319148936104 320.5133333333333,8.021276595744638 320.95,8.946808510638293 321.38666666666666,9.563829787233916 321.8233333333333,9.563829787233916 322.26,9.872340425531947 322.69666666666666,9.563829787233916 323.1333333333333,9.255319148936104 323.57,12.957446808510722 324.00666666666666,12.648936170212691 324.4433333333333,12.031914893617067 324.88,10.7978723404256 325.31666666666666,11.723404255319036 325.75333333333333,12.34042553191488 326.19,12.031914893617067 326.62666666666667,12.031914893617067 327.06333333333333,11.723404255319036 327.5,11.106382978723413 327.93666666666667,9.872340425531947 328.37333333333333,9.563829787233916 328.81,8.946808510638293 329.24666666666667,8.946808510638293 329.68333333333334,7.712765957446827 330.12,8.021276595744638 330.5566666666667,8.021276595744638 330.99333333333334,7.712765957446827 331.43,7.404255319148796 331.8666666666667,6.787234042553172 332.30333333333334,6.787234042553172 332.74,6.478723404255361 333.1766666666667,6.17021276595733 333.61333333333334,7.712765957446827 334.05,7.404255319148796 334.4866666666667,7.712765957446827 334.92333333333335,8.946808510638293 335.36,9.255319148936104 335.7966666666666,8.63829787234048 336.2333333333333,8.63829787234048 336.66999999999996,9.872340425531947 337.1066666666666,9.563829787233916 337.5433333333333,9.255319148936104 337.97999999999996,10.489361702127571 338.41666666666663,10.7978723404256 338.8533333333333,12.648936170212691 339.28999999999996,11.723404255319036 339.72666666666663,12.957446808510722 340.1633333333333,12.34042553191488 340.59999999999997,12.031914893617067 341.03666666666663,10.489361702127571 341.4733333333333,11.106382978723413 341.90999999999997,10.489361702127571 342.34666666666664,11.414893617021225 342.7833333333333,11.723404255319036 343.21999999999997,11.723404255319036 343.65666666666664,11.414893617021225 344.0933333333333,11.106382978723413 344.53,10.180851063829758 344.96666666666664,9.563829787233916 345.4033333333333,9.255319148936104 345.84,6.478723404255361 346.27666666666664,8.021276595744638 346.7133333333333,8.63829787234048 347.15,8.63829787234048 347.58666666666664,8.021276595744638 348.0233333333333,8.32978723404245 348.46,8.946808510638293 348.89666666666665,8.946808510638293 349.3333333333333,9.872340425531947 349.77,15.425531914893654 350.20666666666665,12.957446808510722 350.6433333333333,16.042553191489276 351.08,12.34042553191488 351.51666666666665,14.191489361702187 351.9533333333333,12.34042553191488 352.39,12.957446808510722 352.82666666666665,13.265957446808534 353.2633333333333,11.414893617021225 353.7,12.031914893617067 354.13666666666666,13.882978723404158 354.5733333333333,13.882978723404158 355.01,13.882978723404158 355.44666666666666,13.882978723404158 355.8833333333333,13.882978723404158 356.32,13.574468085106345 356.75666666666666,12.957446808510722 357.1933333333333,13.265957446808534 357.63,12.648936170212691 358.06666666666666,12.648936170212691 358.50333333333333,12.031914893617067 358.94,12.34042553191488 359.37666666666667,12.031914893617067 359.81333333333333,12.648936170212691 360.25,12.648936170212691 360.68666666666667,12.957446808510722 361.12333333333333,12.648936170212691 361.56,12.34042553191488 361.99666666666667,11.723404255319036 362.43333333333334,11.723404255319036 362.87,9.872340425531947 363.3066666666667,10.489361702127571 363.74333333333334,10.180851063829758 364.18,10.7978723404256 364.6166666666667,10.7978723404256 365.05333333333334,10.7978723404256 365.49,10.489361702127571 365.9266666666667,10.7978723404256 366.36333333333334,10.180851063829758 366.8,10.180851063829758 367.2366666666667,10.489361702127571 367.6733333333333,11.106382978723413 368.10999999999996,12.031914893617067 368.5466666666666,10.7978723404256 368.9833333333333,10.7978723404256 369.41999999999996,10.489361702127571 369.8566666666666,11.723404255319036 370.2933333333333,11.723404255319036 370.72999999999996,12.34042553191488 371.16666666666663,12.648936170212691 371.6033333333333,12.34042553191488 372.03999999999996,12.957446808510722 372.47666666666663,13.882978723404158 372.9133333333333,13.574468085106345 373.34999999999997,13.882978723404158 373.78666666666663,13.574468085106345 374.2233333333333,14.191489361702187 374.65999999999997,14.191489361702187 375.09666666666664,14.5 375.5333333333333,16.351063829787307 375.96999999999997,16.351063829787307 376.40666666666664,17.276595744680744 376.8433333333333,17.276595744680744 377.28,18.81914893617024 377.71666666666664,18.51063829787243 378.1533333333333,18.81914893617024 378.59,19.12765957446805 379.02666666666664,19.12765957446805 379.4633333333333,18.51063829787243 379.9,18.81914893617024 380.33666666666664,19.12765957446805 380.7733333333333,18.51063829787243 381.21,19.12765957446805 381.64666666666665,19.744680851063894 382.0833333333333,20.97872340425536 382.52,20.053191489361705 382.95666666666665,29 383.3933333333333,25.914893617021225 383.83,23.13829787234048 384.26666666666665,22.521276595744638 384.7033333333333,23.44680851063829 385.14,21.904255319149016 385.57666666666665,22.212765957446827 386.0133333333333,21.904255319149016 386.45,22.82978723404245 386.88666666666666,22.82978723404245 387.3233333333333,23.44680851063829 387.76,23.44680851063829 388.19666666666666,23.755319148936103 388.6333333333333,24.68085106382976 389.07,24.372340425531945 389.50666666666666,23.755319148936103 389.9433333333333,24.063829787234134 390.38,24.063829787234134 390.81666666666666,24.68085106382976 391.25333333333333,24.68085106382976 391.69,24.063829787234134 392.12666666666667,24.063829787234134 392.56333333333333,24.063829787234134 393,23.755319148936103 393.43666666666667,22.82978723404245 393.87333333333333,22.212765957446827 394.31,22.212765957446827 394.74666666666667,22.212765957446827 395.18333333333334,22.82978723404245 395.62,22.82978723404245 396.0566666666667,21.287234042553173 396.49333333333334,21.595744680850984 396.93,21.595744680850984 397.3666666666667,21.287234042553173 397.80333333333334,21.904255319149016 398.24,21.287234042553173 398.6766666666667,22.521276595744638 399.11333333333334,22.521276595744638 399.54999999999995,21.595744680850984 399.9866666666666,21.904255319149016 400.4233333333333,20.97872340425536 400.85999999999996,21.595744680850984 401.2966666666666,21.287234042553173 401.7333333333333,20.670212765957547 402.16999999999996,21.904255319149016 402.6066666666666,21.287234042553173 403.0433333333333,21.904255319149016 403.47999999999996,21.595744680850984 403.91666666666663,22.521276595744638 404.3533333333333,21.904255319149016 404.78999999999996,21.904255319149016 405.22666666666663,21.287234042553173 405.6633333333333,21.287234042553173 406.09999999999997,20.670212765957547 406.53666666666663,21.595744680850984 406.9733333333333,20.670212765957547 407.40999999999997,20.97872340425536 407.84666666666664,20.361702127659516 408.2833333333333,20.361702127659516 408.71999999999997,20.670212765957547 409.15666666666664,19.744680851063894 409.5933333333333,20.053191489361705 410.03,19.12765957446805 410.46666666666664,19.744680851063894 410.9033333333333,20.053191489361705 411.34,20.053191489361705 411.77666666666664,20.053191489361705 412.2133333333333,19.436170212765862 412.65,18.51063829787243 413.08666666666664,19.12765957446805 413.5233333333333,19.12765957446805 413.96,19.12765957446805 414.39666666666665,18.51063829787243 414.8333333333333,19.436170212765862 415.27,19.12765957446805 415.70666666666665,19.436170212765862 416.1433333333333,19.436170212765862 416.58,20.361702127659516 417.01666666666665,20.670212765957547 417.4533333333333,20.361702127659516 417.89,20.97872340425536 418.32666666666665,20.97872340425536 418.7633333333333,22.82978723404245 419.2,24.063829787234134 419.63666666666666,20.97872340425536 420.0733333333333,22.82978723404245 420.51,25.606382978723413 420.94666666666666,22.82978723404245 421.3833333333333,21.595744680850984 421.82,20.97872340425536 422.25666666666666,20.670212765957547 422.6933333333333,20.361702127659516 423.13,20.97872340425536 423.56666666666666,21.595744680850984 424.00333333333333,20.053191489361705 424.44,20.361702127659516 424.87666666666667,20.670212765957547 425.31333333333333,20.670212765957547 425.75,19.744680851063894 426.18666666666667,21.595744680850984 426.62333333333333,20.97872340425536 427.06,19.744680851063894 427.49666666666667,20.361702127659516 427.93333333333334,20.97872340425536 428.37,23.44680851063829 428.8066666666667,22.82978723404245 429.24333333333334,22.82978723404245 429.68,22.521276595744638 430.1166666666667,17.893617021276587 430.55333333333334,18.51063829787243 430.99,18.51063829787243 431.4266666666666,18.81914893617024 431.8633333333333,17.893617021276587 432.29999999999995,16.351063829787307 432.7366666666666,16.65957446808512 433.1733333333333,14.808510638297811 433.60999999999996,16.65957446808512 434.0466666666666,16.042553191489276 434.4833333333333,15.425531914893654 434.91999999999996,13.574468085106345 435.3566666666666,14.5 435.7933333333333,14.191489361702187 436.22999999999996,14.191489361702187 436.66666666666663,12.34042553191488 437.1033333333333,12.031914893617067 437.53999999999996,11.414893617021225 437.97666666666663,9.563829787233916 438.4133333333333,9.255319148936104 438.84999999999997,8.946808510638293 439.28666666666663,10.7978723404256 439.7233333333333,10.489361702127571 440.15999999999997,10.7978723404256 440.59666666666664,10.7978723404256 441.0333333333333,11.106382978723413 441.46999999999997,10.7978723404256 441.90666666666664,10.180851063829758 442.3433333333333,10.180851063829758 442.78,9.563829787233916 443.21666666666664,10.7978723404256 443.6533333333333,9.563829787233916 444.09,8.63829787234048 444.52666666666664,8.946808510638293 444.9633333333333,9.255319148936104 445.4,8.021276595744638 445.83666666666664,7.404255319148796 446.2733333333333,7.712765957446827 446.71,6.478723404255361 447.14666666666665,8.32978723404245 447.5833333333333,7.404255319148796 448.02,8.63829787234048 448.45666666666665,9.255319148936104 448.8933333333333,8.32978723404245 449.33,6.787234042553172 449.76666666666665,7.712765957446827 450.2033333333333,9.255319148936104 450.64,8.32978723404245 451.07666666666665,11.414893617021225 451.5133333333333,10.180851063829758 451.95,12.031914893617067 452.38666666666666,10.489361702127571 452.8233333333333,10.489361702127571 453.26,11.723404255319036 453.69666666666666,12.031914893617067 454.1333333333333,11.723404255319036 454.57,12.031914893617067 455.00666666666666,11.106382978723413 455.4433333333333,12.648936170212691 455.88,11.414893617021225 456.31666666666666,12.34042553191488 456.75333333333333,13.265957446808534 457.19,14.5 457.62666666666667,14.191489361702187 458.06333333333333,13.574468085106345 458.5,11.723404255319036 458.93666666666667,11.723404255319036 459.37333333333333,12.031914893617067 459.81,10.7978723404256 460.24666666666667,10.180851063829758 460.68333333333334,9.872340425531947 461.12,9.255319148936104 461.5566666666667,9.563829787233916 461.99333333333334,7.712765957446827 462.43,8.63829787234048 462.8666666666667,8.946808510638293 463.3033333333333,9.563829787233916 463.73999999999995,9.872340425531947 464.1766666666666,10.180851063829758 464.6133333333333,8.946808510638293 465.04999999999995,8.63829787234048 465.4866666666666,7.404255319148796 465.9233333333333,8.32978723404245 466.35999999999996,8.021276595744638 466.7966666666666,8.32978723404245 467.2333333333333,7.404255319148796 467.66999999999996,7.404255319148796 468.1066666666666,6.478723404255361 468.5433333333333,6.17021276595733 468.97999999999996,8.32978723404245 469.41666666666663,8.021276595744638 469.8533333333333,8.021276595744638 470.28999999999996,8.021276595744638 470.72666666666663,6.478723404255361 471.1633333333333,7.712765957446827 471.59999999999997,8.021276595744638 472.03666666666663,9.563829787233916 472.4733333333333,8.946808510638293 472.90999999999997,7.404255319148796 473.34666666666664,8.021276595744638 473.7833333333333,7.095744680850984 474.21999999999997,7.404255319148796 474.65666666666664,7.712765957446827 475.0933333333333,8.63829787234048 475.53,8.63829787234048 475.96666666666664,8.32978723404245 476.4033333333333,9.563829787233916 476.84,10.7978723404256 477.27666666666664,11.106382978723413 477.7133333333333,11.723404255319036 478.15,10.7978723404256 478.58666666666664,9.872340425531947 479.0233333333333,10.489361702127571 479.46,11.106382978723413 479.89666666666665,11.414893617021225 480.3333333333333,10.489361702127571 480.77,10.180851063829758 481.20666666666665,10.489361702127571 481.6433333333333,10.7978723404256 482.08,10.489361702127571 482.51666666666665,10.7978723404256 482.9533333333333,11.414893617021225 483.39,11.106382978723413 483.82666666666665,11.106382978723413 484.2633333333333,10.7978723404256 484.7,11.414893617021225 485.13666666666666,11.414893617021225 485.5733333333333,11.414893617021225 486.01,11.723404255319036 486.44666666666666,12.031914893617067 486.8833333333333,12.34042553191488 487.32,12.031914893617067 487.75666666666666,12.957446808510722 488.1933333333333,12.957446808510722 488.63,13.265957446808534 489.06666666666666,12.957446808510722 489.50333333333333,12.648936170212691 489.94,12.957446808510722 490.37666666666667,12.957446808510722 490.81333333333333,13.265957446808534 491.25,12.957446808510722 491.68666666666667,12.648936170212691 492.12333333333333,12.957446808510722 492.56,12.648936170212691 492.99666666666667,12.031914893617067 493.43333333333334,13.574468085106345 493.87,13.574468085106345 494.3066666666667,14.191489361702187 494.74333333333334,13.265957446808534 495.17999999999995,13.574468085106345 495.6166666666666,14.5 496.0533333333333,13.882978723404158 496.48999999999995,13.574468085106345 496.9266666666666,13.574468085106345 497.3633333333333,13.265957446808534 497.79999999999995,12.34042553191488 498.2366666666666,13.882978723404158 498.6733333333333,12.957446808510722 499.10999999999996,12.031914893617067 499.5466666666666,14.191489361702187 499.9833333333333,13.882978723404158 500.41999999999996,13.882978723404158 500.8566666666666,13.882978723404158 501.2933333333333,13.574468085106345 501.72999999999996,13.574468085106345 502.16666666666663,12.957446808510722 502.6033333333333,12.957446808510722 503.03999999999996,13.574468085106345 503.47666666666663,12.648936170212691 503.9133333333333,12.031914893617067 504.34999999999997,11.106382978723413 504.78666666666663,10.7978723404256 505.2233333333333,12.34042553191488 505.65999999999997,12.34042553191488 506.09666666666664,12.031914893617067 506.5333333333333,12.031914893617067 506.96999999999997,10.489361702127571 507.40666666666664,11.414893617021225 507.8433333333333,12.031914893617067 508.28,12.34042553191488 508.71666666666664,12.34042553191488 509.1533333333333,12.34042553191488 509.59,12.031914893617067 510.02666666666664,11.414893617021225 510.4633333333333,11.106382978723413 510.9,12.031914893617067 511.33666666666664,12.031914893617067 511.7733333333333,12.957446808510722 512.2099999999999,13.574468085106345 512.6466666666666,12.34042553191488 513.0833333333333,12.957446808510722 513.52,13.265957446808534 513.9566666666666,12.957446808510722 514.3933333333333,12.957446808510722 514.8299999999999,12.957446808510722 515.2666666666667,13.574468085106345 515.7033333333333,13.574468085106345 516.14,13.574468085106345 516.5766666666666,14.5 517.0133333333333,13.882978723404158 517.4499999999999,13.574468085106345 517.8866666666667,14.5 518.3233333333333,14.191489361702187 518.76,12.34042553191488 519.1966666666666,13.882978723404158 519.6333333333333,13.882978723404158 520.0699999999999,14.5 520.5066666666667,14.5 520.9433333333333,13.882978723404158 521.38,13.882978723404158 521.8166666666666,13.882978723404158 522.2533333333333,12.957446808510722 522.6899999999999,14.808510638297811 523.1266666666667,13.574468085106345 523.5633333333333,13.574468085106345 524,13.574468085106345" style="fill: none;stroke: #416df9;" stroke-dasharray="" stroke-width="undefined"></polyline><text x="0" y="3" fill="#50525c" dominant-baseline="text-before-edge" text-anchor="start" style="font-size: 12px; font-weight: 400; font-family: ;"></text>,<text x="0" y="26" fill="#50525c" dominant-baseline="" text-anchor="start" style="font-size: 12px; font-weight: 400; font-family: ;"></text></svg><!--380--></div>
                        
                        <div class="thumbnail-panel-view" data-s-9dda213a="" style="background: url(&quot;data:image/svg+xml;base64,PHN2ZyB2ZXJzaW9uPSIxLjEiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIKICAgICAgICAgICAgd2lkdGg9IjUyNHB4IgogICAgICAgICAgICBoZWlnaHQ9IjQ3cHgiCiAgICAgICAgPgogICAgICAgICAgICA8cmVjdCB4PSI0MjguMjAwMDAwMDAwMDAwMDUiIHk9IjAiIHdpZHRoPSI4Ni43OTk5OTk5OTk5OTk5NSIgaGVpZ2h0PSIzMCIgc3R5bGU9ImZpbGw6ICM0MTZkZjk7IG9wYWNpdHk6IDAuMTU7IGN1cnNvcjogbW92ZTsiLz4KICAgICAgICAgICAgPGRlZnM+CiAgICAgICAgICAgICAgICA8ZmlsdGVyIGlkPSJ0aHVtYm5haWxQYW5lbFNoYWRvdyI+CiAgICAgICAgICAgICAgICAgICAgPGZlRHJvcFNoYWRvdyBkeD0iMC41IiBkeT0iMC41IiBzdGREZXZpYXRpb249IjEiIGZsb29kLWNvbG9yPSJyZ2JhKDAsIDAsIDAsIDAuMikiIC8+CiAgICAgICAgICAgICAgICAgICAgPGZlRHJvcFNoYWRvdyBkeD0iLTAuNSIgZHk9Ii0wLjUiIHN0ZERldmlhdGlvbj0iMSIgZmxvb2QtY29sb3I9InJnYmEoMCwgMCwgMCwgMC4yKSIgLz4KICAgICAgICAgICAgICAgIDwvZmlsdGVyPgogICAgICAgICAgICA8L2RlZnM+CiAgICAgICAgCiAgICAgICAgICAgIDxyZWN0IHg9IjQxOS4yMDAwMDAwMDAwMDAwNSIgeT0iMSIgd2lkdGg9IjE4IiBoZWlnaHQ9IjI4IiByeD0iNiIgcnk9IjYiIGZpbGw9IiNmZmYiIGZpbHRlcj0idXJsKCN0aHVtYm5haWxQYW5lbFNoYWRvdykiLz4KICAgICAgICAgICAgPHJlY3QgeD0iNDIzLjIwMDAwMDAwMDAwMDA1IiB5PSIxMCIgd2lkdGg9IjIiIGhlaWdodD0iMTAiIHJ4PSIxIiByeT0iMSIgZmlsbD0iI0UwRTNGMCIvPgogICAgICAgICAgICA8cmVjdCB4PSI0MjcuMjAwMDAwMDAwMDAwMDUiIHk9IjEyIiB3aWR0aD0iMiIgaGVpZ2h0PSI2IiByeD0iMSIgcnk9IjEiIGZpbGw9IiNFMEUzRjAiLz4KICAgICAgICAgICAgPHJlY3QgeD0iNDMxLjIwMDAwMDAwMDAwMDA1IiB5PSIxMCIgd2lkdGg9IjIiIGhlaWdodD0iMTAiIHJ4PSIxIiByeT0iMSIgZmlsbD0iI0UwRTNGMCIvPgoKICAgICAgICAgICAgPHJlY3QgeD0iNTA2IiB5PSIxIiB3aWR0aD0iMTgiIGhlaWdodD0iMjgiIHJ4PSI2IiByeT0iNiIgZmlsbD0iI2ZmZiIgZmlsdGVyPSJ1cmwoI3RodW1ibmFpbFBhbmVsU2hhZG93KSIvPgogICAgICAgICAgICA8cmVjdCB4PSI1MTAiIHk9IjEwIiB3aWR0aD0iMiIgaGVpZ2h0PSIxMCIgcng9IjEiIHJ5PSIxIiBmaWxsPSIjRTBFM0YwIi8+CiAgICAgICAgICAgIDxyZWN0IHg9IjUxNCIgeT0iMTIiIHdpZHRoPSIyIiBoZWlnaHQ9IjYiIHJ4PSIxIiByeT0iMSIgZmlsbD0iI0UwRTNGMCIvPgogICAgICAgICAgICA8cmVjdCB4PSI1MTgiIHk9IjEwIiB3aWR0aD0iMiIgaGVpZ2h0PSIxMCIgcng9IjEiIHJ5PSIxIiBmaWxsPSIjRTBFM0YwIi8+CiAgICAgICAgPC9zdmc+&quot;) center center no-repeat;"></div><!--381-->
                        <div class="feedback" data-s-9dda213a=""></div>
                    <!--379-->
    <object tabindex="-1" type="text/html" aria-hidden="true" data="about:blank" style="display: block; position: absolute; top: 0px; left: 0px; width: 100%; height: 100%; border: none; padding: 0px; margin: 0px; opacity: 0; z-index: -1000; pointer-events: none;"></object></div>
                </div><!--377-->
            </div>

            <div class="trade-data-container" data-s-9dda213a="" style="height: 389px;">
                <div class="five-container" data-s-90ba8f78="" style="width: 237px;">
        <div class="flod-bar" data-s-90ba8f78="">
            <i class="cos-icon cos-icon-right"></i>
        </div>
        <div class="five" data-s-90ba8f78="">
            <!--387-->
                <div class="five_title" data-s-90ba8f78="">
                    <span class="info" data-s-90ba8f78="">买
                    <div class="bar" data-s-90ba8f78="">
                        <div class="buy item" data-s-90ba8f78="" style="width: 24.6%;"></div>
                        <div class="sell item" data-s-90ba8f78="" style="width: 75.4%;"></div>
                    </div>
                    <span class="info" data-s-90ba8f78="">卖
                </div>
                <ul class="five_box" data-s-90ba8f78="">
                    <li data-s-90ba8f78="">
                        <div class="item " data-s-90ba8f78="">
                            <div class="text" data-s-90ba8f78="">卖5</div>
                            <div class="price" data-s-90ba8f78="" style="color: var(--stock-color-red);">36.19</div>
                            <div class="volume" data-s-90ba8f78="">
                                60
                                
                                <div class="volumebg green-bg" data-s-90ba8f78="" style="width: 23.7%;"></div>
                            </div>
                            
                            <div class="mask  green-mask" data-s-90ba8f78=""></div>
                        </div>
                        <!--389-->
                    </li><li data-s-90ba8f78="">
                        <div class="item " data-s-90ba8f78="">
                            <div class="text" data-s-90ba8f78="">卖4</div>
                            <div class="price" data-s-90ba8f78="" style="color: var(--stock-color-red);">36.18</div>
                            <div class="volume" data-s-90ba8f78="">
                                54
                                
                                <div class="volumebg green-bg" data-s-90ba8f78="" style="width: 21.3%;"></div>
                            </div>
                            
                            <div class="mask  green-mask" data-s-90ba8f78=""></div>
                        </div>
                        <!--390-->
                    </li><li data-s-90ba8f78="">
                        <div data-s-90ba8f78="" class="item ">
                            <div class="text" data-s-90ba8f78="">卖3</div>
                            <div class="price" data-s-90ba8f78="" style="color: var(--stock-color-red);">36.17</div>
                            <div class="volume" data-s-90ba8f78="">
                                253
                                
                                <div data-s-90ba8f78="" class="volumebg green-bg" style="width: 100%;"></div>
                            </div>
                            
                            <div data-s-90ba8f78="" class="mask  green-mask"></div>
                        </div>
                        <!--391-->
                    </li><li data-s-90ba8f78="">
                        <div data-s-90ba8f78="" class="item ">
                            <div class="text" data-s-90ba8f78="">卖2</div>
                            <div class="price" data-s-90ba8f78="" style="color: var(--stock-color-red);">36.16</div>
                            <div class="volume" data-s-90ba8f78="">
                                72
                                
                                <div data-s-90ba8f78="" class="volumebg green-bg" style="width: 28.5%;"></div>
                            </div>
                            
                            <div data-s-90ba8f78="" class="mask  green-mask"></div>
                        </div>
                        <!--392-->
                    </li><li data-s-90ba8f78="">
                        <div data-s-90ba8f78="" class="item ">
                            <div class="text" data-s-90ba8f78="">卖1</div>
                            <div class="price" data-s-90ba8f78="" style="color: var(--stock-color-red);">36.14</div>
                            <div class="volume" data-s-90ba8f78="">
                                91
                                
                                <div data-s-90ba8f78="" class="volumebg green-bg" style="width: 35.9%;"></div>
                            </div>
                            
                            <div data-s-90ba8f78="" class="mask  green-mask"></div>
                        </div>
                        <div class="line" data-s-90ba8f78=""></div><!--393-->
                    </li><li data-s-90ba8f78="">
                        <div data-s-90ba8f78="" class="item ">
                            <div class="text" data-s-90ba8f78="">买1</div>
                            <div class="price" data-s-90ba8f78="" style="color: var(--stock-color-red);">36.13</div>
                            <div class="volume" data-s-90ba8f78="">
                                4
                                
                                <div data-s-90ba8f78="" class="volumebg red-bg" style="width: 1.6%;"></div>
                            </div>
                            
                            <div data-s-90ba8f78="" class="mask  red-mask"></div>
                        </div>
                        <!--394-->
                    </li><li data-s-90ba8f78="">
                        <div data-s-90ba8f78="" class="item ">
                            <div class="text" data-s-90ba8f78="">买2</div>
                            <div class="price" data-s-90ba8f78="" style="color: var(--stock-color-red);">36.10</div>
                            <div class="volume" data-s-90ba8f78="">
                                4
                                
                                <div data-s-90ba8f78="" class="volumebg red-bg" style="width: 1.6%;"></div>
                            </div>
                            
                            <div data-s-90ba8f78="" class="mask  red-mask"></div>
                        </div>
                        <!--395-->
                    </li><li data-s-90ba8f78="">
                        <div data-s-90ba8f78="" class="item ">
                            <div class="text" data-s-90ba8f78="">买3</div>
                            <div class="price" data-s-90ba8f78="" style="color: var(--stock-color-red);">36.08</div>
                            <div class="volume" data-s-90ba8f78="">
                                3
                                
                                <div data-s-90ba8f78="" class="volumebg red-bg" style="width: 1.2%;"></div>
                            </div>
                            
                            <div data-s-90ba8f78="" class="mask  red-mask"></div>
                        </div>
                        <!--396-->
                    </li><li data-s-90ba8f78="">
                        <div data-s-90ba8f78="" class="item ">
                            <div class="text" data-s-90ba8f78="">买4</div>
                            <div class="price" data-s-90ba8f78="" style="color: var(--stock-color-red);">36.07</div>
                            <div class="volume" data-s-90ba8f78="">
                                27
                                
                                <div data-s-90ba8f78="" class="volumebg red-bg" style="width: 10.7%;"></div>
                            </div>
                            
                            <div data-s-90ba8f78="" class="mask  red-mask"></div>
                        </div>
                        <!--397-->
                    </li><li data-s-90ba8f78="">
                        <div data-s-90ba8f78="" class="item ">
                            <div class="text" data-s-90ba8f78="">买5</div>
                            <div class="price" data-s-90ba8f78="" style="color: var(--stock-color-red);">36.06</div>
                            <div class="volume" data-s-90ba8f78="">
                                135
                                
                                <div data-s-90ba8f78="" class="volumebg red-bg" style="width: 53.4%;"></div>
                            </div>
                            
                            <div data-s-90ba8f78="" class="mask  red-mask"></div>
                        </div>
                        <!--398-->
                    </li><!--388-->
                </ul>
                <p class="gap-line" data-s-90ba8f78="">
                    <span class="text" data-s-90ba8f78="">分笔成交
                    <i class="cos-icon cos-icon-up"></i>
                </p>
            <!--387--><!--386-->
            <div class="detail-wrapper" data-s-90ba8f78="">
                <ul data-s-90ba8f78="">
                    <li class="item" data-s-90ba8f78="">
                        <span class="mask green-mask " data-s-90ba8f78="">
                        <span class="text" data-s-90ba8f78="">15:00
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            778
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--401-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span class="mask red-mask " data-s-90ba8f78="">
                        <span class="text" data-s-90ba8f78="">14:57
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            23
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--402-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:56
                        <span class="price" data-s-90ba8f78="">36.10
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--403-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:56
                        <span class="price" data-s-90ba8f78="">36.15
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--404-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:56
                        <span class="price" data-s-90ba8f78="">36.10
                        <span class="volume" data-s-90ba8f78="">
                            12
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--405-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:56
                        <span class="price" data-s-90ba8f78="">36.12
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--406-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:56
                        <span class="price" data-s-90ba8f78="">36.07
                        <span class="volume" data-s-90ba8f78="">
                            9
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--407-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:56
                        <span class="price" data-s-90ba8f78="">36.06
                        <span class="volume" data-s-90ba8f78="">
                            8
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--408-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:56
                        <span class="price" data-s-90ba8f78="">36.06
                        <span class="volume" data-s-90ba8f78="">
                            11
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--409-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:56
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--410-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:56
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--411-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:56
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            5
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--412-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:56
                        <span class="price" data-s-90ba8f78="">36.08
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--413-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:56
                        <span class="price" data-s-90ba8f78="">36.08
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--414-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:56
                        <span class="price" data-s-90ba8f78="">36.06
                        <span class="volume" data-s-90ba8f78="">
                            8
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--415-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:56
                        <span class="price" data-s-90ba8f78="">36.08
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--416-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:56
                        <span class="price" data-s-90ba8f78="">36.12
                        <span class="volume" data-s-90ba8f78="">
                            291
                            <span class="normal-gray" data-s-90ba8f78="" style="color: var(--stock-color-flat);">M<!--417-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:56
                        <span class="price" data-s-90ba8f78="">36.12
                        <span class="volume" data-s-90ba8f78="">
                            58
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--418-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:55
                        <span class="price" data-s-90ba8f78="">36.16
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--419-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:55
                        <span class="price" data-s-90ba8f78="">36.15
                        <span class="volume" data-s-90ba8f78="">
                            6
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--420-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:55
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            5
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--421-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:55
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            15
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--422-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:55
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            4
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--423-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:55
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--424-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:55
                        <span class="price" data-s-90ba8f78="">36.15
                        <span class="volume" data-s-90ba8f78="">
                            20
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--425-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:55
                        <span class="price" data-s-90ba8f78="">36.15
                        <span class="volume" data-s-90ba8f78="">
                            4
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--426-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:55
                        <span class="price" data-s-90ba8f78="">36.12
                        <span class="volume" data-s-90ba8f78="">
                            8
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--427-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:55
                        <span class="price" data-s-90ba8f78="">36.15
                        <span class="volume" data-s-90ba8f78="">
                            16
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--428-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:55
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--429-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:55
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            6
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--430-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:55
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            43
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--431-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:55
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--432-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:55
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--433-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:55
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--434-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:55
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            5
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--435-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:55
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--436-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:55
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            5
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--437-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:54
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--438-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:54
                        <span class="price" data-s-90ba8f78="">36.15
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--439-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:54
                        <span class="price" data-s-90ba8f78="">36.15
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--440-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:54
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--441-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:54
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            63
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--442-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:54
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--443-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:54
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            7
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--444-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:54
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            25
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--445-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:54
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--446-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:54
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            6
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--447-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:54
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            7
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--448-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:54
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            4
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--449-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:54
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--450-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:54
                        <span class="price" data-s-90ba8f78="">36.15
                        <span class="volume" data-s-90ba8f78="">
                            6
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--451-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:54
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            8
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--452-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:54
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            21
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--453-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:54
                        <span class="price" data-s-90ba8f78="">36.15
                        <span class="volume" data-s-90ba8f78="">
                            4
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--454-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:54
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--455-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:54
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            57
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--456-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:53
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            6
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--457-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:53
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--458-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:53
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            4
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--459-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:53
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            18
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--460-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:53
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            9
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--461-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:53
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            4
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--462-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:53
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--463-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:53
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            4
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--464-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:53
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            4
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--465-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:53
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            14
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--466-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:53
                        <span class="price" data-s-90ba8f78="">36.15
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--467-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:53
                        <span class="price" data-s-90ba8f78="">36.12
                        <span class="volume" data-s-90ba8f78="">
                            90
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--468-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:53
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--469-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:53
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--470-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:53
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--471-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:52
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--472-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:52
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            8
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--473-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:52
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            16
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--474-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:52
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--475-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:52
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--476-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:52
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--477-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:52
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            4
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--478-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:52
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            4
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--479-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:52
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--480-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:52
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            50
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--481-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:52
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--482-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:52
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            5
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--483-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:52
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--484-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:52
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            20
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--485-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:52
                        <span class="price" data-s-90ba8f78="">36.12
                        <span class="volume" data-s-90ba8f78="">
                            12
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--486-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:52
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--487-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:51
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--488-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:51
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--489-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:51
                        <span class="price" data-s-90ba8f78="">36.10
                        <span class="volume" data-s-90ba8f78="">
                            7
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--490-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:51
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--491-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:51
                        <span class="price" data-s-90ba8f78="">36.10
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--492-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:51
                        <span class="price" data-s-90ba8f78="">36.10
                        <span class="volume" data-s-90ba8f78="">
                            5
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--493-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:51
                        <span class="price" data-s-90ba8f78="">36.12
                        <span class="volume" data-s-90ba8f78="">
                            10
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--494-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:51
                        <span class="price" data-s-90ba8f78="">36.12
                        <span class="volume" data-s-90ba8f78="">
                            9
                            <span class="normal-gray" data-s-90ba8f78="" style="color: var(--stock-color-flat);">M<!--495-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:51
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            55
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--496-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:51
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--497-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:51
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--498-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:51
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--499-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:51
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--500-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:51
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            10
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--501-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:51
                        <span class="price" data-s-90ba8f78="">36.12
                        <span class="volume" data-s-90ba8f78="">
                            5
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--502-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:50
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            9
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--503-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:50
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            7
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--504-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:50
                        <span class="price" data-s-90ba8f78="">36.12
                        <span class="volume" data-s-90ba8f78="">
                            14
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--505-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:50
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            5
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--506-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:50
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            4
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--507-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:50
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--508-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:50
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--509-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:50
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--510-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:50
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            4
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--511-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:50
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            100
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--512-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:50
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            4
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--513-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:50
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--514-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:50
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            17
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--515-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:50
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            132
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--516-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:50
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            6
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--517-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:50
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            18
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--518-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:50
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            108
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--519-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:49
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            69
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--520-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:49
                        <span class="price" data-s-90ba8f78="">36.12
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--521-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:49
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            69
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--522-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:49
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--523-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:49
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            49
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--524-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:49
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--525-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:49
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            134
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--526-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:49
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            7
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--527-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:49
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--528-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:49
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            28
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--529-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:49
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--530-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:49
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            4
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--531-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:49
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            21
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--532-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:49
                        <span class="price" data-s-90ba8f78="">36.12
                        <span class="volume" data-s-90ba8f78="">
                            21
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--533-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:49
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            100
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--534-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:49
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            69
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--535-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:49
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            6
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--536-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:49
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            4
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--537-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:49
                        <span class="price" data-s-90ba8f78="">36.16
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--538-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:48
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--539-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:48
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--540-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:48
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            18
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--541-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:48
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            5
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--542-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:48
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            4
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--543-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:48
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--544-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:48
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="normal-gray" data-s-90ba8f78="" style="color: var(--stock-color-flat);">M<!--545-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:48
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--546-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:48
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--547-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:48
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--548-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:48
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            16
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--549-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:48
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--550-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:48
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            4
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--551-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:48
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            69
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--552-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:48
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--553-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:48
                        <span class="price" data-s-90ba8f78="">36.17
                        <span class="volume" data-s-90ba8f78="">
                            47
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--554-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:48
                        <span class="price" data-s-90ba8f78="">36.15
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--555-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:47
                        <span class="price" data-s-90ba8f78="">36.18
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--556-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:47
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--557-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:47
                        <span class="price" data-s-90ba8f78="">36.18
                        <span class="volume" data-s-90ba8f78="">
                            108
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--558-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:47
                        <span class="price" data-s-90ba8f78="">36.16
                        <span class="volume" data-s-90ba8f78="">
                            38
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--559-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:47
                        <span class="price" data-s-90ba8f78="">36.15
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--560-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:47
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            16
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--561-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:47
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            1
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--562-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:47
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--563-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:47
                        <span class="price" data-s-90ba8f78="">36.15
                        <span class="volume" data-s-90ba8f78="">
                            7
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--564-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:47
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            23
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--565-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:47
                        <span class="price" data-s-90ba8f78="">36.16
                        <span class="volume" data-s-90ba8f78="">
                            6
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--566-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:47
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            7
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--567-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:47
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            22
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--568-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:47
                        <span class="price" data-s-90ba8f78="">36.16
                        <span class="volume" data-s-90ba8f78="">
                            65
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--569-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:47
                        <span class="price" data-s-90ba8f78="">36.15
                        <span class="volume" data-s-90ba8f78="">
                            14
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--570-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:46
                        <span class="price" data-s-90ba8f78="">36.12
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--571-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:46
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--572-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:46
                        <span class="price" data-s-90ba8f78="">36.12
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--573-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:46
                        <span class="price" data-s-90ba8f78="">36.12
                        <span class="volume" data-s-90ba8f78="">
                            13
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--574-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:46
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            9
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--575-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:46
                        <span class="price" data-s-90ba8f78="">36.12
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--576-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:46
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--577-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:46
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--578-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:46
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--579-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:46
                        <span class="price" data-s-90ba8f78="">36.13
                        <span class="volume" data-s-90ba8f78="">
                            9
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--580-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:46
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--581-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:46
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--582-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:46
                        <span class="price" data-s-90ba8f78="">36.15
                        <span class="volume" data-s-90ba8f78="">
                            3
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--583-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:46
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            12
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--584-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:46
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            4
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--585-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:46
                        <span class="price" data-s-90ba8f78="">36.15
                        <span class="volume" data-s-90ba8f78="">
                            6
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--586-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:46
                        <span class="price" data-s-90ba8f78="">36.15
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--587-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:45
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            4
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--588-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:45
                        <span class="price" data-s-90ba8f78="">36.11
                        <span class="volume" data-s-90ba8f78="">
                            8
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--589-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:45
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--590-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:45
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            6
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--591-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:45
                        <span class="price" data-s-90ba8f78="">36.10
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--592-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:45
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            5
                            <span class="normal-gray" data-s-90ba8f78="" style="color: var(--stock-color-flat);">M<!--593-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:45
                        <span class="price" data-s-90ba8f78="">36.10
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--594-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:45
                        <span class="price" data-s-90ba8f78="">36.15
                        <span class="volume" data-s-90ba8f78="">
                            11
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--595-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:45
                        <span class="price" data-s-90ba8f78="">36.10
                        <span class="volume" data-s-90ba8f78="">
                            2
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--596-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask red-mask ">
                        <span class="text" data-s-90ba8f78="">14:45
                        <span class="price" data-s-90ba8f78="">36.15
                        <span class="volume" data-s-90ba8f78="">
                            43
                            <span class="up bs" data-s-90ba8f78="" style="color: var(--stock-color-red);">B<!--597-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:45
                        <span class="price" data-s-90ba8f78="">36.10
                        <span class="volume" data-s-90ba8f78="">
                            8
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--598-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:45
                        <span class="price" data-s-90ba8f78="">36.10
                        <span class="volume" data-s-90ba8f78="">
                            10
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--599-->
                        
                    </li><li class="item" data-s-90ba8f78="">
                        <span data-s-90ba8f78="" class="mask green-mask ">
                        <span class="text" data-s-90ba8f78="">14:45
                        <span class="price" data-s-90ba8f78="">36.14
                        <span class="volume" data-s-90ba8f78="">
                            16
                            <span class="down bs" data-s-90ba8f78="" style="color: var(--stock-color-green);">S<!--600-->
                        
                    </li><!--400-->
                </ul>
            </div>
        </div><!--385-->
    </div>
            </div><!--382-->
        </div>
    </div><!--344-->
            </div>
        

                        <div class="price harmony-os-bold" data-s-3dd71b5f="">36.14</div>
                        <div class="unit cos-font-medium" data-s-3dd71b5f="">元</div>
                        <div class="increase harmony-os-medium" data-s-3dd71b5f="">+0.30</div>
                        <div class="ratio harmony-os-medium" data-s-3dd71b5f="">+0.84%</div>
                    

                                
                                
                                <svg style="transform: rotate(-180deg); margin-bottom: 2px;" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="8px" height="6px" viewBox="0 0 12 8" version="1.1" data-s-d84eef82="">
                                    <g id="控件" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" data-s-d84eef82="">
                                        <g id="盘口数据/默认展示3行" transform="translate(-181.000000, -130.000000)" fill="#858585" data-s-d84eef82="">
                                            <g id="编组" transform="translate(24.000000, 88.000000)" data-s-d84eef82="">
                                                <g id="拖动/横滑组件备份-2" transform="translate(118.000000, 26.000000)" data-s-d84eef82="">
                                                    <g id="编组-2" transform="translate(39.000000, 16.000000)" data-s-d84eef82="">
                                                        
                                                        <path d="M6.8479983,1.35679729 L10.0437507,6.47000106 C10.336461,6.93833759 10.1940878,7.55528797 9.7257513,7.8479983 C9.5668197,7.94733056 9.38317208,8 9.19575236,8 L2.80424764,8 C2.25196289,8 1.80424764,7.55228475 1.80424764,7 C1.80424764,6.81258028 1.85691709,6.62893266 1.95624934,6.47000106 L5.1520017,1.35679729 C5.44471203,0.888460755 6.06166241,0.74608759 6.52999894,1.03879792 C6.65876807,1.11927863 6.7675176,1.22802815 6.8479983,1.35679729 Z" id="三角形备份" transform="translate(6.000000, 4.000000) scale(1, -1) rotate(-180.000000) translate(-6.000000, -4.000000) " data-s-d84eef82=""></path>
                                                    </g>
                                                </g>
                                            </g>
                                        </g>
                                    </g>
                                </svg>
                            <!--886-->
            
        <!--885-->
    
    <!--314-->
                         
                            <!--331-->
                         
                            <!--326-->
                         
                            <!--330-->
                         0 36.14`

	_ = s
}
func getByTagBaidu(src, tag, end string) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	ss := strings.Split(src, tag)[1]
	ss = strings.Split(src, `font-weight: 500;">`)[1]
	if end != "" {
		return strings.Split(ss, end)[0]
	}
	return strings.Split(ss, `</span>`)[0]
}
func getDetailBaidu(html, shtml string, s Stock) Stock {
	fmt.Println(shtml)
	//html = `<div class="stock-info"><div class="stock-price stock-fall"><div class="stock-current"><strong>¥18.03</strong></div><div class="stock-change">-0.12  -0.66%</div></div><div class="stock-time"><div>&nbsp;53.55 万球友关注</div><div class="quote-market-status"><span>交易中<pan> 09-30 14:12:03 北京时间</span></div></div></div><table class="quote-info"><tbody><tr><td>最高：<span class="stock-rise">18.20</spa<td>今开：<span class="stock-fall">18.09</span></td><td>涨停：<span class="stock-rise">19.97</span></td><td>成交量：<span>66.24万手</sp><tr class="separateTop"><td>最低：<span class="stock-fall">17.71</span></td><td>昨收：<span>18.15</span></td><td>跌停：<span class="st>16.34</span></td><td>成交额：<span>11.85亿</span></td></tr><tr class="separateBottom"><td>量比：<span class="stock-fall">0.66</span></手：<span>0.34%</span></td><td>市盈率(动)：<span>9.95</span></td><td>市盈率(TTM)：<span>10.66</span></td></tr><tr><td>委比：<span class-40.71%</span></td><td>振幅：<span>2.70%</span></td><td>市盈率(静)：<span>12.10</span></td><td>市净率：<span>1.14</span></td></tr><tr><n>1.69</span></td><td>股息(TTM)：<span>0.18</span></td><td>总股本：<span>194.06亿</span></td><td>总市值：<span>3498.89亿</span></td></t资产：<span>15.83</span></td><td>股息率(TTM)：<span>1.00%</span></td><td>流通股：<span>194.06亿</span></td><td>流通值：<span>3498.86亿<tr><td>52周最高：<span>25.16</span></td><td>52周最低：<span>14.64</span></td><td>货币单位：<span>CNY</span></td></tr></tbody></table>`
	zsz := strings.ReplaceAll(getByTagBaidu(html, `总市值`, "</div>"), "亿", "")
	ltz := getByTagBaidu(html, `流通值`, "</div>")
	ttm := getByTagBaidu(html, `市盈(TTM)`, "</div>")
	sylj := getByTagBaidu(html, `市盈(静)`, "</div>")
	sum := getByTagBaidu(html, `总股本`, "</div>")
	shtm := getByTagBaidu(shtml, `price harmony-os-bold`, "</div>")
	price := getByTagBaidu(shtm, ">", "")
	Zsz, err := strconv.ParseFloat(zsz, 64)
	if err != nil {
		Zsz = 0
	}
	s.Zsz = Zsz

	//re := regexp.MustCompile("[\u4e00-\u9fa5]{1,}")
	//service := re.FindAllString(strings.ReplaceAll(getByTag(shtml, `<div class="title">业务</div>`, `<!---->`), ",", "，"), -1) //业务
	//name := re.FindAllString(strings.ReplaceAll(getByTag(shtml, `<h1 class="stock-name">`, `(`), ",", "，"), -1)              //name
	//info := re.FindAllString(strings.ReplaceAll(getByTag(shtml, `<div class="title">简介</div>`, `<!---->`), ",", "，"), -1)    //简介

	//s.Name = strings.Join(name, "，")
	//s.Service = strings.Join(service, "，")
	//s.Info = strings.Join(info, "，")

	s.Ltz = ltz
	Ttm, err := strconv.ParseFloat(ttm, 64)
	if err != nil {
		Ttm = 0
	}
	s.Ttm = Ttm
	Sylj, err := strconv.ParseFloat(sylj, 64)
	if err != nil {
		Sylj = 0
	}
	s.Sylj = Sylj
	c := 1.0
	if strings.Contains(sum, "万") {
		sum = strings.ReplaceAll(sum, "万", "")
		c = 10000
	} else if strings.Contains(sum, "亿") {
		c = 100000000
		sum = strings.ReplaceAll(sum, "亿", "")
	}
	Sum, err := strconv.ParseFloat(sum, 64)
	if err != nil {
		Sum = 0.0
	}
	s.Sum = Sum * c
	Price, err := strconv.ParseFloat(price, 64)
	if err != nil {
		Price = 0.0
	}
	s.Price = Price
	//sy := getByTag(html,`市盈`)
	fmt.Println(zsz, ltz, ttm, sylj, s.Sum, s.Price)
	return s
}

type Stock struct {
	Name    string  `名称`
	Service string  `业务`
	Info    string  `简介`
	code    string  `代码`
	Jzcsyl  float64 `净资产收益率`
	Mll     float64 `毛利率`
	Zsz     float64 `总市值`
	Ltz     string  `流通市值`
	Ttm     float64 `TTM`
	Sylj    float64 `市盈率(静)`
	Sum     float64 `总股本`
	Price   float64 `当前价格`
}
