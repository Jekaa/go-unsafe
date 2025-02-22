package main

import (
	"fmt"
	"runtime"
	"unsafe"
)

// Структура с приватными и публичными полями
type myStruct struct {
	PublicField  int
	privateField int
}

// Пакетная структура с приватными полями
type privateStruct struct {
	hiddenField int
}

func main() {
	// 1. Базовые преобразования типов
	demoBasicConversions()

	// 2. Доступ к приватным полям
	demoPrivateAccess()

	// 3. Арифметика указателей
	demoPointerArithmetic()

	// 4. Опасности uintptr и GC
	demoUintptrDanger()
}

func demoBasicConversions() {
	fmt.Println("=== Basic Conversions ===")

	// Преобразование между типами через unsafe.Pointer
	var i int = 42
	fmt.Printf("Original int: %d\n", i)

	// int -> float64
	f := *(*float64)(unsafe.Pointer(&i))
	fmt.Printf("As float64: %f\n", f)

	// Обратное преобразование
	f = 3.1415
	i = *(*int)(unsafe.Pointer(&f))
	fmt.Printf("Float as int: %d\n", i)
	fmt.Println()
}

func demoPrivateAccess() {
	fmt.Println("=== Private Field Access ===")

	// Создание экземпляра структуры
	s := myStruct{PublicField: 100, privateField: 200}
	fmt.Printf("Before: Public=%d, private=%d\n", s.PublicField, s.privateField)

	// Получение доступа к приватному полю
	ptr := unsafe.Pointer(&s)
	privateAddr := uintptr(ptr) + unsafe.Offsetof(s.privateField)

	// Изменение приватного поля
	*(*int)(unsafe.Pointer(privateAddr)) = 300
	fmt.Printf("After: Public=%d, private=%d\n", s.PublicField, s.privateField)

	// Доступ к полям пакетной структуры
	ps := privateStruct{hiddenField: 42}
	psPtr := unsafe.Pointer(&ps)
	hidden := (*int)(unsafe.Pointer(uintptr(psPtr) + unsafe.Offsetof(ps.hiddenField)))
	fmt.Printf("Accessed hidden field: %d\n", *hidden)
	fmt.Println()
}

func demoPointerArithmetic() {
	fmt.Println("=== Pointer Arithmetic ===")

	arr := [3]int{10, 20, 30}
	ptr := unsafe.Pointer(&arr[0])

	for i := 0; i < 3; i++ {
		// Вычисление адреса элемента массива
		addr := uintptr(ptr) + uintptr(i)*unsafe.Sizeof(arr[0])
		val := *(*int)(unsafe.Pointer(addr))
		fmt.Printf("arr[%d] = %d\n", i, val)
	}

	// Работа со структурой
	s := struct {
		a byte
		b int
		c float64
	}{}

	size := unsafe.Sizeof(s)
	align := unsafe.Alignof(s.b)
	offset := unsafe.Offsetof(s.c)

	fmt.Printf("Size: %d, Align(b): %d, Offset(c): %d\n", size, align, offset)
	fmt.Println()
}

func demoUintptrDanger() {
	fmt.Println("=== uintptr GC Danger ===")

	// Создаем объект
	create := func() uintptr {
		data := []int{1, 2, 3}
		ptr := unsafe.Pointer(&data[0])
		return uintptr(ptr)
	}

	addr := create()

	// Вызываем сборщик мусора
	runtime.GC()

	// Попытка использования адреса после GC (опасно!)
	data := *(*int)(unsafe.Pointer(addr))
	fmt.Printf("Value at invalid address: %d (может быть мусором)\n", data)
	fmt.Println("Это демонстрация опасности - результат непредсказуем!")
}
