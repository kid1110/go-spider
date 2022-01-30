package controller

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/solenovex/web-tutor/exception"
	"github.com/solenovex/web-tutor/model"
	"github.com/solenovex/web-tutor/result"
	"github.com/solenovex/web-tutor/service"
)

func RegisterNewsRoutes() {
	http.HandleFunc("/getNews", exception.ErrWrapper(FlushHandler))
	http.HandleFunc("/searchByTitle", exception.ErrWrapper(SearchByTitleHandler))
	http.HandleFunc("/SearchByDate", exception.ErrWrapper(SearchByDateHandler))
}

func FlushHandler(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	confirt := r.Form.Get("confirt")
	if confirt != "1" {
		return exception.CreateNewException(result.ParameterErrorException, result.GetStautsString(result.ParameterErrorException))
	}
	//爬取并且存入
	url := service.GetUrl()
	err := service.GetNews(url)
	if err != nil {
		return err
	}
	//清空redis库存
	service.MyRedis.RedisClear()

	w.Header().Set("Content-Type", "application/json")
	var res = result.GetStatusText(result.Success)
	err = json.NewEncoder(w).Encode(&res)
	if err != nil {
		log.Fatal("return error!")
	}
	return nil

}

func SearchByTitleHandler(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	title := r.Form.Get("title")
	if title == "" {
		return exception.CreateNewException(result.ParameterErrorException, result.GetStautsString(result.ParameterErrorException))
	}
	//从redis中获取数据
	n := []model.News{}
	err := service.MyRedis.RedisGet("title:"+title, &n)
	if err != nil {
		return err
	}
	//当redis有库存，拿redis的库存
	if len(n) != 0 {
		var res = result.GetStatusText(result.Success)
		res.Data = n
		err := json.NewEncoder(w).Encode(&res)
		if err != nil {
			log.Fatal("return error!")
			return err
		}
		return nil

	} else {
		loc := sync.Mutex{}
		//否则查询mysql
		loc.Lock()
		data := service.FindNewByTitle(title)
		loc.Unlock()
		if len(data) != 0 {
			//存入redis
			rand.Seed(time.Now().Unix())
			num := rand.Intn(10)
			//设置随机时间
			service.MyRedis.RedisAdd("title:"+title, data, time.Duration(num)*time.Minute)
			var res = result.GetStatusText(result.Success)
			res.Data = data
			err = json.NewEncoder(w).Encode(&res)
			if err != nil {
				log.Fatal("return error!")
				return err
			}
		} else {
			//返回无信息
			res := result.GetStatusText(result.NoTitleData)
			err := json.NewEncoder(w).Encode(&res)
			if err != nil {
				log.Fatal("return error!")
				return err
			}

		}
		return nil
	}

}
func SearchByDateHandler(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	date := r.Form.Get("date")
	//先从布隆过滤器找
	judge := service.MyBloom.Exist(date)
	//如果没找到这说明是穿透，直接返回找不到
	if !judge {
		var res = result.GetStatusText(result.Success)
		res.Data = "暂无符合当天日期的新闻"
		err := json.NewEncoder(w).Encode(&res)
		if err != nil {
			log.Fatal("return error!")
		}
		return nil

	}
	//再在redis中找
	n := []model.News{}

	err := service.MyRedis.RedisGet("date:"+date, &n)
	if err != nil {
		return err
	}
	if len(n) != 0 {
		var res = result.GetStatusText(result.Success)
		res.Data = n
		err := json.NewEncoder(w).Encode(&res)
		if err != nil {
			log.Fatal("return error!")
			return err
		}
		return nil
	} else {
		lock := sync.Mutex{}
		//需要在mysql中找
		lock.Lock()
		data := service.GetNewsByDate(date)
		lock.Unlock()
		if len(data) != 0 {
			//需要存redis中
			rand.Seed(time.Now().Unix())
			num := rand.Intn(10)
			//设置随机时间
			service.MyRedis.RedisAdd("date:"+date, data, time.Duration(num)*time.Minute)

			//返回需要的值
			var res = result.GetStatusText(result.Success)
			res.Data = data
			err = json.NewEncoder(w).Encode(&res)
			if err != nil {
				log.Fatal("return error!")
				return err
			}
			return nil
		}

	}

	return nil
}
