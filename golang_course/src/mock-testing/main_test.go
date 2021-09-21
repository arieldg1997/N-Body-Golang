package main

import "testing"

func TestGetFullTimeExployeeById(t *testing.T) {
	table := []struct {
		id               int
		dni              string
		mockFunc         func()
		expectedExployee FullTimeEmployee
	}{
		{
			id:  1,
			dni: "1",
			mockFunc: func() {
				GetEmployeeById = func(id int) (Employee, error) {
					return Employee{
						Id:       1,
						Position: "CEO",
					}, nil
				}

				GetPersonByDNI = func(id string) (Person, error) {
					return Person{
						Name: "Ariel",
						Age:  23,
						DNI:  "1",
					}, nil
				}
			},
			expectedExployee: FullTimeEmployee{
				Person: Person{
					Age:  23,
					DNI:  "1",
					Name: "Ariel",
				},
				Employee: Employee{
					Id:       1,
					Position: "CEO",
				},
			},
		},
	}
	originalGetEmployeeById := GetEmployeeById
	originalGetPersonByDNI := GetPersonByDNI
	for _, test := range table {
		test.mockFunc()
		ft, err := GetFullTimeEmployeeById(test.id, test.dni)
		if err != nil {
			t.Errorf("error when getting employee")
		}
		if ft.Age != test.expectedExployee.Age {
			t.Errorf("error, got %d expected %d", ft.Age, test.expectedExployee.Age)
		}
		// Asi con todas las demas propiedades...
	}
	GetEmployeeById = originalGetEmployeeById
	GetPersonByDNI = originalGetPersonByDNI
}
