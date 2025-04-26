package models

type PersonFilter struct {
    Name        string 
    Surname     string 
    Age         int   
    Gender      string 
    Nationality string 
    Page  int
    Limit int 

}

// Offset возвращает смещение для SQL LIMIT/OFFSET
func (f PersonFilter) Offset() int {
    if f.Page <= 1 {
        return 0
    }
    return (f.Page - 1) * f.Limit
}