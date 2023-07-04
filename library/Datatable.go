package library

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//ykzirla
//DV8uJq2UWP6BgNa
//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

type DataTableConfig struct {
	Primary    string
	Where      string
	Group      string
	Order      string
	Result     interface{}
	Column     []map[string]interface{}
	Table      string
	Connection *gorm.DB
	Request    map[string]interface{}
	TableCount string
}
type DataTableOutput struct {
	Draw            int         `json:"draw"`
	RecordsTotal    int         `json:"recordsTotal"`
	RecordsFiltered int         `json:"recordsFiltered"`
	Data            interface{} `json:"data"`
}

func DataTableTest() string {
	return "Hello Success"
}

func DataTable(Config DataTableConfig) (interface{}, error) {
	dbReader := Config.Connection

	var queryBuilder string

	strColumn := DTGenerateColumn(Config.Column)
	queryBuilder = queryBuilder + strColumn + " FROM " + Config.Table + " "

	//building where
	var whereDefault, dtWhere string
	if len(Config.Where) > 0 {
		whereDefault = " WHERE " + Config.Where
	}
	dtWhere = DTGenerateWhere(Config.Column, Config.Request)
	if len(dtWhere) > 0 {
		if len(whereDefault) > 0 {
			dtWhere = " AND (" + dtWhere + ")"
		} else {
			dtWhere = "WHERE " + dtWhere
		}
	}

	queryBuilder = queryBuilder + whereDefault + dtWhere
	//end building where
	var defaultGroup string
	if len(Config.Group) > 0 {
		defaultGroup = " GROUP BY " + Config.Group
	}
	queryBuilder = queryBuilder + defaultGroup

	//build order
	var orderDefault, dtOrder string

	if len(Config.Order) > 0 {
		orderDefault = " ORDER BY " + Config.Order
	}
	dtOrder = DTGenerateOrder(Config.Column, Config.Request)
	if len(dtOrder) > 0 {
		if len(orderDefault) > 0 {
			dtOrder = ", " + dtOrder
		} else {
			dtOrder = " ORDER BY " + dtOrder
		}
	}
	queryBuilder = queryBuilder + orderDefault + dtOrder

	//end building order
	total := 0
	if len(Config.TableCount) > 0 {
		readCount := dbReader.Raw("SELECT COUNT(" + Config.Primary + ") FROM " + Config.TableCount + " " + whereDefault + " " + defaultGroup + " LIMIT 2").Scan(&total)
		if readCount.RowsAffected > 1 {
			dbReader.Raw("SELECT SUM(a.total) FROM (SELECT COUNT(" + Config.Primary + ") as total FROM " + Config.TableCount + " " + whereDefault + " " + defaultGroup + ") as a").Scan(&total)
		}
	} else {
		readCount := dbReader.Raw("SELECT COUNT(" + Config.Primary + ") FROM " + Config.Table + " " + whereDefault + " " + defaultGroup + " LIMIT 2").Scan(&total)
		if readCount.RowsAffected > 1 {
			dbReader.Raw("SELECT SUM(a.total) FROM (SELECT COUNT(" + Config.Primary + ") as total FROM " + Config.Table + " " + whereDefault + " " + defaultGroup + ") as a").Scan(&total)
		}
		fmt.Println(readCount.RowsAffected)
	}
	fmt.Println(total)

	filteredtotal := total
	// filteredCount := dbReader.Raw("SELECT COUNT(" + Config.Primary + ") FROM " + Config.Table + " " + whereDefault + dtWhere + " " + defaultGroup).Scan(&filteredtotal)
	// if filteredCount.RowsAffected > 1 {
	// 	dbReader.Raw("SELECT SUM(a.total) FROM (SELECT COUNT(" + Config.Primary + ") as total FROM " + Config.Table + " " + whereDefault + dtWhere + " " + defaultGroup + ") as a").Scan(&filteredtotal)
	// }

	//limit builder
	queryBuilder = queryBuilder + DTGeneratorLimit(Config.Request)
	//end limit builder
	rows, err := dbReader.Raw("SELECT " + queryBuilder).Rows()
	if err != nil {
		return nil, err
	}
	//string column
	defer rows.Close()
	var dataOutput []map[string]interface{}
	//var output []map[string]interface{}
	for rows.Next() {
		colOuputstring := "{"
		for _, v := range Config.Column {
			colOuputstring = colOuputstring + fmt.Sprintf("%v", v["alias"]) + ","
		}
		colOuputstring = colOuputstring + "}"
		j, _ := JsonDecode(colOuputstring)
		dbReader.ScanRows(rows, &j)
		dataOutput = append(dataOutput, j)
	}
	if len(dataOutput) == 0 {
		dataOutput = make([]map[string]interface{}, 0)
	}
	outputConcept := map[string]interface{}{"draw": 1, "recordsTotal": total, "recordsFiltered": filteredtotal, "data": dataOutput}
	return outputConcept, nil
}

func DTGenerateColumn(col []map[string]interface{}) string {
	var strColumn string
	strColumn = ""
	for i, s := range col {
		if i == 0 {
			strColumn = fmt.Sprintf("%v", s["col"]) + " as " + fmt.Sprintf("%v", s["alias"])
		} else {
			strColumn = strColumn + ", " + fmt.Sprintf("%v", s["col"]) + " as " + fmt.Sprintf("%v", s["alias"])
		}
	}
	return strColumn
}

func DTGenerateWhere(col []map[string]interface{}, req map[string]interface{}) string {
	if req["cols"] != nil && req["search"] != nil {
		searchValue := fmt.Sprintf("%v", req["search"])
		if len(searchValue) > 0 {

			colsJ, _ := JsonEncode(req["cols"])
			cols, _ := JsonDecode(colsJ)

			whereOuput := ""
			for _, v := range col {
				if cols[fmt.Sprintf("%v", v["alias"])] != nil {
					searchJ, _ := JsonEncode(cols[fmt.Sprintf("%v", v["alias"])])
					search, _ := JsonDecode(searchJ)
					var colName string
					if search["searchable"].(bool) {
						colName = fmt.Sprintf("%v", v["col"])
						if fmt.Sprintf("%v", v["type"]) == "number" {
							colName = "CAST(" + colName + " as TEXT)"
						}
						if len(whereOuput) == 0 {
							whereOuput = colName + " ILIKE '%" + searchValue + "%'"
						} else {
							whereOuput = whereOuput + " OR " + colName + " ILIKE '%" + searchValue + "%'"
						}
					}
				}
			}
			return whereOuput
		}
	}
	return ""
}

func DTGenerateOrder(col []map[string]interface{}, req map[string]interface{}) string {
	var orderOuput string
	if req["orderBy"] != nil && req["orderType"] != nil {
		orderBy := fmt.Sprintf("%v", req["orderBy"])
		orderType := fmt.Sprintf("%v", req["orderType"])
		for _, v := range col {
			if len(orderBy) > 0 && len(orderType) > 0 && v["alias"] == orderBy {
				orderOuput = orderOuput + fmt.Sprintf("%v", v["col"]) + " " + orderType
			}
		}
	}
	return orderOuput
}

func DTGeneratorLimit(req map[string]interface{}) string {
	page := 1
	if req["page"] != nil {
		pageInterface := fmt.Sprintf("%v", req["page"])
		p, err := strconv.Atoi(pageInterface)
		if err != nil {
			fmt.Println(err.Error())
		}
		page = p
	}

	nRecords := 10
	if req["nRecords"] != nil {
		pageNRecord := fmt.Sprintf("%v", req["nRecords"])
		p, err := strconv.Atoi(pageNRecord)
		if err != nil {
			fmt.Println(err.Error())
		}
		nRecords = p
	}

	var limitOutput string
	var startLimit, lengthLimit int
	startLimit = page*nRecords - nRecords
	lengthLimit = nRecords

	limitOutput = " LIMIT " + fmt.Sprintf("%v", lengthLimit) + " OFFSET " + fmt.Sprintf("%v", startLimit)
	return limitOutput
}

func DTExportToExcel(Config DataTableConfig) (string, error) {
	dbReader := Config.Connection

	var queryBuilder string
	//queryBuilder = "SELECT "

	//column Name
	strColumn := DTGenerateColumn(Config.Column)
	queryBuilder = queryBuilder + strColumn + " FROM " + Config.Table + " "

	//building where
	var whereDefault, dtWhere string
	if len(Config.Where) > 0 {
		whereDefault = " WHERE " + Config.Where
	}

	dtWhere = DTGenerateWhere(Config.Column, Config.Request)
	if len(dtWhere) > 0 {
		if len(whereDefault) > 0 {
			dtWhere = " AND " + dtWhere
		} else {
			dtWhere = "WHERE " + dtWhere
		}
	}
	queryBuilder = queryBuilder + whereDefault + dtWhere
	//end building where
	var defaultGroup string
	if len(Config.Group) > 0 {
		defaultGroup = " GROUP BY " + Config.Group
	}
	queryBuilder = queryBuilder + defaultGroup
	//build order
	var orderDefault, dtOrder string

	if len(Config.Order) > 0 {
		orderDefault = " ORDER BY " + Config.Order
	}
	dtOrder = DTGenerateOrder(Config.Column, Config.Request)
	if len(dtOrder) > 0 {
		if len(orderDefault) > 0 {
			dtOrder = ", " + dtOrder
		} else {
			dtOrder = " ORDER BY " + dtOrder
		}
	}
	queryBuilder = queryBuilder + orderDefault + dtOrder
	//end building order
	xlsx := excelize.NewFile()
	sheet1Name := "Sheet 1"
	xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)
	inc := 0
	//header
	var mapResult []string
	for _, v := range Config.Column {
		export := v["export"].(bool)
		if export {
			xlsx.SetCellValue(sheet1Name, string('A'+inc)+"1", strings.ToUpper(fmt.Sprintf("%v", v["alias"])))
			mapResult = append(mapResult, fmt.Sprintf("%v", v["alias"]))
			inc++
		}
	}

	//data
	rows, err := dbReader.Raw("SELECT " + queryBuilder).Rows()
	if err != nil {
		return "", err
	}
	jsonRes, _ := JsonEncode(mapResult)
	jsonRess, _ := JsonDecode(jsonRes)
	iRows := 2
	for rows.Next() {
		dbReader.ScanRows(rows, &jsonRess)
		//fmt.Println(jsonRess["name"])
		inc = 0
		for _, v := range Config.Column {
			export := v["export"].(bool)
			if export {
				xlsx.SetCellValue(sheet1Name, string('A'+inc)+strconv.Itoa(iRows), fmt.Sprintf("%v", jsonRess[fmt.Sprintf("%v", v["alias"])]))
				inc++
			}
		}
		iRows++
	}
	fName := uuid.New().String() + ".xlsx"
	err = xlsx.SaveAs(os.Getenv("upload_path") + fName)
	if err != nil {
		return "", err
	}
	return fName, nil
}
