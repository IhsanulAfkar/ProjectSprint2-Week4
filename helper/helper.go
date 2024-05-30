package helper

import (
	"math"
	"time"

	"github.com/jmoiron/sqlx"
)

func CheckUsername(conn *sqlx.DB, username string) bool {
	query := "SELECT EXISTS (SELECT 1 FROM public.user WHERE username = $1 UNION SELECT 1 FROM admin WHERE username = $2) AS username_exists"
	var isExist bool
	err := conn.QueryRow(query,username, username).Scan(&isExist)
	if err!= nil{
		 
		return false
	}
	return isExist
}
func FormatToIso860(s string)string {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		
		return ""
	}

	// Format the time object into ISO 8601 format
	return t.Format("2006-01-02T15:04:05Z07:00")
}
func Includes(target string, array []string)bool{
	for _, value := range array {
        if value == target {
            return true
        }
    }
    return false
}

func CountHaversine(lat1 float64, lon1 float64, lat2 float64, lon2 float64)float64{
	Radius := 6371
    lat1 = lat1 * math.Pi / 180
    lon1 = lon1 * math.Pi / 180
    lat2 = lat2 * math.Pi / 180
    lon2 = lon2 * math.Pi / 180
    latDiff := lat2 - lat1
    lonDiff := lon2 - lon1
    a := math.Sin(latDiff/2)*math.Sin(latDiff/2) +
        math.Cos(lat1)*math.Cos(lat2)*math.Sin(lonDiff/2)*math.Sin(lonDiff/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    km := float64(Radius) * c

    return km
}
