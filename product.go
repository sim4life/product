package product

import (
	"bufio"
	// "encoding/json"
	//"flag"
	"fmt"
	//elastigo "github.com/mattbaird/elastigo/lib"
	//"log"
	// "os"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

type PrintStructVals interface {
	PrintValues()
}

type Products []Product

type Product struct {
	ID                        int        `json:"-"`
	Product_ID                int64      `json:"productId"`
	Name                      string     `json:"name"`
	Item_Code                 string     `json:"itemCode"`
	Description               string     `json:"description"`
	Record_State              int64      `json:"-"`
	Used_In_Recipe_ID         []int64    `json:"-"`
	Dep_Comm                  string     `json:"-"`
	Groceries_Category_ID     []int64    `json:"-"`
	Groceries_Categories      [][]string `json:"groceriesCategories"`
	Groceries_Search_Synonyms []string   `json:"gorceriesSearchSynonyms"`
	Ideas_Category_ID         []int64    `json:"-"`
	Ideas_Categories          [][]string `json:"ideasCategories"`
	Ideas_Search_Synonyms     []string   `json:"ideasSearchSynonyms"`
}

type SliceField struct {
	field  string
	values []int64
	isOld  bool
}

func (products Products) PrintValues() {
	for ind, prod := range products {
		fmt.Printf("product[%d]:\n%+v\n", ind, prod)
	}
}

func (products Products) PrintJSON() {
	if productsJSON, err := json.MarshalIndent(products, "", "   "); err == nil {
		fmt.Println("JSON is: ", string(productsJSON))
	}
}

func (products Products) FlattenGroceriesCategories(cat_map MapCategories) {
	fmt.Println("In Products ==> FlattenGroceriesCategories\n\n")
	fmt.Printf("len(products) is: %d and len(cat_map) is: %d\n", len(products), len(cat_map))

	for i, prod := range products {
		for _, groc_cat_id := range prod.Groceries_Category_ID {
			flat_groc_cat := &FlatCategory{Category_ID: groc_cat_id}
			flat_groc_cat.FlattenCategory(cat_map)

			products[i].Groceries_Categories = append(prod.Groceries_Categories, flat_groc_cat.Category_Names)
			products[i].Groceries_Search_Synonyms = append(prod.Groceries_Search_Synonyms, flat_groc_cat.Search_Synonyms)
		}
		// fmt.Printf("---product is:%+v\n", prod)
		// fmt.Printf("prod.Groceries_Categories is:%q\n", prod.Groceries_Categories)

	}

	// fmt.Printf("produsts outside loop is: %+v\n", products)
	// products.PrintValues()
	fmt.Println("Ending FlattenGroceriesCategories---------")

	// products.PrintValues()

	// return products
}

func (products Products) FlattenIdeasCategories(cat_map MapCategories) {
	fmt.Println("In Products ==> FlattenIdeasCategories\n\n")
	fmt.Printf("len(products) is: %d and len(cat_map) is: %d\n", len(products), len(cat_map))
}

func setCurrentSliceField(product *Product, sliceField *SliceField) {
	// fmt.Printf("-/-/--sliceField.isOld is:%v\n", sliceField.isOld)
	if sliceField.isOld {
		// fmt.Printf("str setting values:%v\n", sliceField.values)
		reflect.ValueOf(product).Elem().FieldByName(sliceField.field).Set(reflect.ValueOf(sliceField.values))
		sliceField.isOld = false
	}

}

func readSliceField(sliceField *SliceField, product *Product, field string, value int64) SliceField {
	// fmt.Printf("sField:%s while field:%s == %v\n", sliceField.field, field, sliceField.isOld)
	// fmt.Printf("the condition is: %v\n", (sliceField.field != field || !sliceField.isOld))

	// if !sliceField.isOld {
	if sliceField.field != field || !sliceField.isOld {
		setCurrentSliceField(product, sliceField)

		sliceField = &SliceField{field, []int64{value}, true}
		// fmt.Printf("  --Inner sliceField is: %+v\n", sliceField)
	} else {
		sliceField.values = append(sliceField.values, value)
	}

	// fmt.Printf(" -----sliceField is: %+v\n", sliceField)
	return *sliceField
}

func readSetField(product *Product, field, value string, setField func(field reflect.Value, value string)) string {
	fieldName := reflect.ValueOf(product).Elem().FieldByName(field)
	if fieldName.IsValid() {
		// fieldName.SetInt(int64(value))
		setField(fieldName, value)
		return ""
	} else {
		return field
	}

}

func setStrIntField(sliceField *SliceField, unknowns []string, product *Product, field, value string,
	setField func(field reflect.Value, value string)) {

	setCurrentSliceField(product, sliceField)

	unknown := readSetField(product, field, value, setField)
	if unknown != "" {
		unknowns = append(unknowns, unknown)
	}

}

func ParseProductFile(scanner *bufio.Scanner) PrintStructVals {
	counter := 1
	product := Product{ID: counter}
	var products Products
	var sliceField SliceField
	unknowns := make([]string, 0)

	for scanner.Scan() {
		str := strings.Split(scanner.Text(), "|")
		// fmt.Printf("str[0]:%s ", str[0])
		if str[0] != "EOR" {
			// fmt.Printf("str[1]:%s\n", str[1])

			field_kind := reflect.ValueOf(product).FieldByName(str[0]).Kind()
			if field_kind == reflect.String || field_kind == reflect.Invalid {
				setStrIntField(&sliceField, unknowns, &product, str[0], str[1],
					func(field reflect.Value, value string) {
						field.SetString(value)
					})
			} else {
				if value, err := strconv.ParseInt(str[1], 10, 64); err == nil {
					//TODO: get the value of field from the struct and then manipulate it
					// instead of relying on data input as contiguous
					if field_kind == reflect.Slice {
						sliceField = readSliceField(&sliceField, &product, str[0], value)
						// fmt.Printf("values len is:%d", len(values))
					} else {
						setStrIntField(&sliceField, unknowns, &product, str[0], str[1],
							func(field reflect.Value, value string) {
								int_val, _ := strconv.ParseInt(value, 10, 64)
								field.SetInt(int_val)
							})
					}
				}
			}

		} else {
			//flush any pending slice data
			setCurrentSliceField(&product, &sliceField)
			/*
				if counter%20 == 0 {
					fmt.Println("\n\n{\"index\" : {\"_id\" :", counter, ", \"_index\":\"policereport\", \"_type\":\"crimestat\"}}")
					fmt.Println("%+v\nFinish----------", product)
				}
			*/
			products = append(products, product)
			counter++
			product = Product{ID: counter}
			// fmt.Println("%+v----------Init", product)
		}
	}

	fmt.Printf("\nSuccessfully read products: %v\n", len(products))
	fmt.Printf("\nNo. of unkowns de-normalized fields are: %v\n", len(unknowns))

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		//TODO: return error
		return nil //, err
	}
	return products

}
