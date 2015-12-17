package product

import (
	"bufio"
	"fmt"
	"reflect"
	// "sort"
	"strconv"
	"strings"
)

type Category struct {
	ID                 int    `json:"-"`
	Category_ID        int64  `json:"categoryId"`
	Parent_Category_ID int64  `json:"-"`
	Category_Name      string `json:"name"`
	Search_Synonyms    string `json:"searchSynonym"`
	Sequence           string `json:"-"`

	parent   *Category   `json:"-"`
	children []*Category `json:"-"`
}

type FlatCategory struct {
	// ID              int
	Category_ID     int64
	Category_Names  []string // each item contains parents name
	Search_Synonyms string   // all parents synonyms appended
}

type MapCategories map[int64]*Category

func (cat_map MapCategories) PrintValues() {
	// var m map[int]string
	var keys []int64
	for k, v := range cat_map {
		keys = append(keys, k)
		fmt.Printf("key: %d and val:%+v\n", k, v)
		fmt.Printf("parent is: %+v\n", v.parent)
		if len(v.children) > 0 {
			fmt.Printf("child[0] is: %+v\n\n", v.children[0])
		}
	}
	fmt.Printf("\nGrandaddy is:%+v\n", cat_map[int64(1)])

}

func (flat_cat *FlatCategory) FlattenCategory(cat_map MapCategories) {

	cat := cat_map[flat_cat.Category_ID]
	for cat != nil {
		flat_cat.Category_Names = append(flat_cat.Category_Names, cat.Category_Name)
		flat_cat.Search_Synonyms += cat.Search_Synonyms + ","
		cat = cat.parent
	}
	flat_cat.Search_Synonyms = strings.Trim(strings.Trim(flat_cat.Search_Synonyms, ","), " ")
	// fmt.Printf("flat_cat is:%+v\n", flat_cat)
	// fmt.Printf("flat_cat.Category_Names is:%q\n", flat_cat.Category_Names)
}

func populateCategoryTree(cat_map MapCategories) {
	counter := 0
	//finding the root of the tree
	var root_groc_cat *Category
	for _, v := range cat_map {
		if v.Parent_Category_ID == int64(0) {
			root_groc_cat = v
			counter++
		}
	}

	// TODO: if counter is greater than 1 then
	// throw error/panic
	fmt.Printf("root is: %+v", root_groc_cat)
	fmt.Println("found roots: ", counter)

	for _, v := range cat_map {
		parent := cat_map[v.Parent_Category_ID]
		if parent != nil {
			parent.children = append(parent.children, v)
			v.parent = parent
		}
	}
}

func ParseCategoriesFile(scanner *bufio.Scanner) PrintStructVals {
	fmt.Println("fun fun")

	counter := 1
	cat_id := int64(0)
	category := &Category{ID: counter}
	cat_map := make(MapCategories)

	for scanner.Scan() {
		str := strings.Split(scanner.Text(), "|")
		if str[0] != "" {
			fmt.Printf("str[0]:%s ", str[0])
			if str[0] != "EOR" {
				fmt.Printf("str[1]:%s\n", str[1])

				field_kind := reflect.ValueOf(*category).FieldByName(str[0]).Kind()
				if field_kind == reflect.String {
					reflect.ValueOf(category).Elem().FieldByName(str[0]).SetString(str[1])
				} else { //field_kind == reflect.Int
					if value, err := strconv.Atoi(str[1]); err == nil {
						if str[0] == "Category_ID" {
							cat_id = int64(value)
						}
						reflect.ValueOf(category).Elem().FieldByName(str[0]).SetInt(int64(value))
					}
				}

			} else {
				// if counter%100 == 0 {
				fmt.Println("{\"index\" : {\"_id\" :", counter, ", \"_index\":\"policereport\", \"_type\":\"crimestat\"}}")
				fmt.Println("%+v\nFinish----------", category)
				// }

				cat_map[cat_id] = category
				counter++
				category = &Category{ID: counter}

				fmt.Println("%+v----------Init", category)

			}
		}
	}

	populateCategoryTree(cat_map)
	// cat_map.PrintValues()

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		//TODO: return proper error
		return nil //err
	}

	return cat_map

}
