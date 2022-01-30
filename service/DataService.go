package service

import (
	"github.com/sirupsen/logrus"
	"github.com/solenovex/web-tutor/model"
)

func CheckGid(gid string) int32 {
	stmt, err := MyDB.Prepare("SELECT COUNT(1) FROM db_news.t_news WHERE gid = ?")
	if err != nil {
		Mylog.WithFields(logrus.Fields{
			"stmt": stmt,
		}).Error("查询gid字段预处理出错:" + err.Error())
	}
	rows, err := stmt.Query(gid)
	defer rows.Close()
	var count int32
	for rows.Next() {
		err = rows.Scan(&count)
	}
	if err != nil {
		Mylog.WithFields(logrus.Fields{
			"rows": rows,
		}).Error("查询gid字段数据出错:" + err.Error())
	}
	return count
}
func CheckDate(date string) int32 {
	stmt, err := MyDB.Prepare("SELECT COUNT(1) FROM db_news.t_news WHERE date = ?")
	if err != nil {
		Mylog.WithFields(logrus.Fields{
			"stmt": stmt,
		}).Error("查询gid字段预处理出错:" + err.Error())
	}
	rows, err := stmt.Query(date)
	defer rows.Close()
	var count int32
	for rows.Next() {
		err = rows.Scan(&count)
	}
	if err != nil {
		Mylog.WithFields(logrus.Fields{
			"rows": rows,
		}).Error("查询date字段数据出错:" + err.Error())
	}
	return count

}
func GetNewsByDate(date string) []model.News {
	stmt, err := MyDB.Prepare("SELECT nid,gid,title,context,date FROM db_news.t_news WHERE date = ?")
	if err != nil {
		Mylog.WithFields(logrus.Fields{
			"stmt": stmt,
		}).Error("使用date查询新闻字段预处理出错:" + err.Error())
	}
	rows, err := stmt.Query(date)
	if err != nil {
		Mylog.WithFields(logrus.Fields{
			"rows": rows,
		}).Error("使用date查询新闻字段出错:" + err.Error())
	}
	defer rows.Close()
	var data []model.News
	for rows.Next() {
		var News = model.News{}
		err = rows.Scan(&News.Nid, &News.Gid, &News.Title, &News.Text, &News.Date)
		if err != nil {
			Mylog.WithFields(logrus.Fields{
				"News": News,
			}).Error("使用date查询新闻字段数据出错:" + err.Error())
		}
		data = append(data, News)
	}
	return data

}

func FindNewByTitle(title string) []model.News {
	stmt, err := MyDB.Prepare("SELECT nid,gid,title,context,date FROM db_news.t_news WHERE title like concat('%',?,'%')")
	if err != nil {
		Mylog.WithFields(logrus.Fields{
			"stmt": stmt,
		}).Error("使用title模糊查询新闻字段预处理出错:" + err.Error())
	}
	rows, err := stmt.Query(title)
	if err != nil {
		Mylog.WithFields(logrus.Fields{
			"rows": rows,
		}).Error("使用title查询新闻字段处理出错:" + err.Error())
	}
	defer rows.Close()
	var data []model.News
	for rows.Next() {
		var News = model.News{}
		err = rows.Scan(&News.Nid, &News.Gid, &News.Title, &News.Text, &News.Date)
		if err != nil {
			Mylog.WithFields(logrus.Fields{
				"News": News,
			}).Error("使用title查询新闻字段处理数据出错:" + err.Error())
		}
		data = append(data, News)
	}
	return data
}
func InsertNews(news model.News) int64 {
	// // tr, err := MyDB.Begin()
	// if err != nil {
	// 	fmt.Println("开启事务失败" + err.Error())
	// }
	stmt, err := MyDB.Prepare("INSERT INTO db_news.t_news( gid, title, context, date) VALUES (?,?,?,?)")
	if err != nil {
		Mylog.WithFields(logrus.Fields{
			"stmt": stmt,
		}).Error("插入新闻预处理出错:" + err.Error())
	}
	defer stmt.Close()
	// defer tr.Commit()

	querry, err := stmt.Exec(news.Gid, news.Title, news.Text, news.Date)
	if err != nil {
		Mylog.WithFields(logrus.Fields{
			"querry": querry,
		}).Error("插入新闻处理出错:" + err.Error())
	}
	result, err := querry.RowsAffected()
	if err != nil {
		Mylog.WithFields(logrus.Fields{
			"result": result,
		}).Error("插入新闻处理数据出错:" + err.Error())
	}
	return result
}
func GetNewByGid(gid string) model.News {
	stmt, err := MyDB.Prepare("SELECT nid,gid,title,context,date FROM db_news.t_news WHERE gid = ? LIMIT 1")
	if err != nil {
		Mylog.WithFields(logrus.Fields{
			"stmt": stmt,
		}).Error("通过gid获取新闻预处理出错:" + err.Error())
	}
	querry, err := stmt.Query(gid)
	defer querry.Close()
	if err != nil {
		Mylog.WithFields(logrus.Fields{
			"querry": querry,
		}).Error("通过gid获取新闻出错:" + err.Error())
	}
	var News = model.News{}
	for querry.Next() {
		err = querry.Scan(&News.Nid, &News.Gid, &News.Title, &News.Text, &News.Date)
		if err != nil {
			Mylog.WithFields(logrus.Fields{
				"News": News,
			}).Error("通过gid获取新闻处理数据出错:" + err.Error())
		}
	}
	return News
}
