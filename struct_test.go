package gobase

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		testFile string
		expected []Schema
	}{
		{
			testFile: "./testdata/create_table.go",
			expected: []Schema{
				{
					SchemaName: "users",
					SchemaFields: []SchemaField{
						{Name: "ID", DataType: "int"},
						{Name: "Name", DataType: "string"},
						{Name: "CreatedAt", DataType: "Time"},
						{Name: "UpdatedAt", DataType: "Time"},
						{Name: "IsMember", DataType: "bool"},
					},
				},
			},
		},
		{
			testFile: "./testdata/multi_table_test.go",
			expected: []Schema{
				{
					SchemaName: "users",
					SchemaFields: []SchemaField{
						{Name: "ID", DataType: "int"},
						{Name: "Name", DataType: "string"},
						{Name: "CreatedAt", DataType: "Time"},
						{Name: "UpdatedAt", DataType: "Time"},
						{Name: "IsMember", DataType: "bool"},
					},
				},
				{
					SchemaName: "image",
					SchemaFields: []SchemaField{
						{Name: "Name", DataType: "string"},
						{Name: "Type", DataType: "string"},
						{Name: "Size", DataType: "int"},
						{Name: "Hidden", DataType: "bool"},
						{Name: "CreatedAt", DataType: "Time"},
						{Name: "UpdatedAt", DataType: "Time"},
					},
				},
			},
		},
		{
			testFile: "./testdata/multi_table_test2.go",
			expected: []Schema{
				{
					SchemaName: "User",
					SchemaFields: []SchemaField{
						{Name: "ID", DataType: "int"},
						{Name: "Name", DataType: "string"},
						{Name: "CreatedAt", DataType: "Time"},
						{Name: "UpdatedAt", DataType: "Time"},
						{Name: "IsMember", DataType: "bool"},
					},
				},
				{
					SchemaName: "Image",
					SchemaFields: []SchemaField{
						{Name: "Name", DataType: "string"},
						{Name: "Type", DataType: "string"},
						{Name: "CreatedAt", DataType: "Time"},
						{Name: "UpdatedAt", DataType: "Time"},
						{Name: "IsMember", DataType: "bool"},
					},
				},
				{
					SchemaName: "Product",
					SchemaFields: []SchemaField{
						{Name: "ID", DataType: "string"},
						{Name: "Name", DataType: "string"},
						{Name: "Description", DataType: "string"},
						{Name: "Price", DataType: "float64"},
						{Name: "InStock", DataType: "bool"},
						{Name: "CreatedAt", DataType: "Time"},
						{Name: "UpdatedAt", DataType: "Time"},
					},
				},
				{
					SchemaName: "Order",
					SchemaFields: []SchemaField{
						{Name: "OrderID", DataType: "string"},
						{Name: "UserID", DataType: "int"},
						{Name: "ProductID", DataType: "string"},
						{Name: "Quantity", DataType: "int"},
						{Name: "TotalPrice", DataType: "float64"},
						{Name: "OrderDate", DataType: "Time"},
						{Name: "Delivered", DataType: "bool"},
					},
				},
				{
					SchemaName: "Review",
					SchemaFields: []SchemaField{
						{Name: "ReviewID", DataType: "string"},
						{Name: "ProductID", DataType: "string"},
						{Name: "UserID", DataType: "int"},
						{Name: "Rating", DataType: "int"},
						{Name: "Comment", DataType: "string"},
						{Name: "ReviewDate", DataType: "Time"},
					},
				},
				{
					SchemaName: "Category",
					SchemaFields: []SchemaField{
						{Name: "CategoryID", DataType: "string"},
						{Name: "Name", DataType: "string"},
						{Name: "ParentID", DataType: "string"},
						{Name: "CreatedAt", DataType: "Time"},
						{Name: "UpdatedAt", DataType: "Time"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		got := Parse(tt.testFile)

		if !reflect.DeepEqual(got, tt.expected) {
			t.Fatalf(
				"Test failed for file: %s\n\nExpected:\n%+v\n\nGot:\n%+v",
				tt.testFile,
				tt.expected,
				got,
			)
		}
	}
}
