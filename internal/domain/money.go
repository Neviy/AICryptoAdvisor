// В пакете domain будут хранится сущности монеты и юзера
// Именно в этом отделе domain хранится	и инициализируется сущность монеты

package domain

import "errors"

// Это структура монеты , которая имеет id для поиска в БД, название , стоимость , процент прибыли или убытка и рекомендация от AI
type Coin struct{
	ID int64               `json:"id"`
	Name string            `json:"name"`
	Price float64          `json:"price"`
	Percent float64        `json:"percent"`
	Recommendation string  `json:"recommendation"`
}
// Конструктор для создания монеты . Сюда подается название и цена монеты , а процент и рекомендация нулевые значения , потому что раньше этих монет не было в системе. 
func NewCoin(name string,price float64)(*Coin,error){
	var errs []error
	if name==""{
		errs = append(errs, errors.New("not correct name"))
	}
	if price <=0{ 
		errs = append(errs, errors.New("not correct price"))}
	if len(errs) >0{
		return nil,errors.Join(errs...)
	}
	return &Coin{
		Name: name,
		Price: price,
		Percent: 0,
		Recommendation: "",
	},nil
}

// Функция для того чтобы притягивать данные из БД . Она возвращает данные без ошибки , ведь данные уже проходили валидацию 
func NewCoinFromDB(id int64,name string,price float64,percent float64,recommendation string)*Coin{
	return &Coin{
		ID: id,
		Name: name,
		Price: price,
		Percent: percent,
		Recommendation: recommendation,
	}
}
// Меняет цену и  пересчитывает  процент
func (c *Coin) UpdatePrice(newPrice float64) error {
	if newPrice <= 0 {
		return errors.New("invalid price")
	}
	if c.Price != 0 {
		c.Percent = ((newPrice - c.Price) / c.Price) * 100
	} else {
		c.Percent = 0
	}
	c.Price = newPrice
	return nil
}