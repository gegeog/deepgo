package main

import (
	"fmt"
	"math"
	"unsafe"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		sl := unsafe.Slice(unsafe.StringData(name), len(name))
		person.name = [42]byte(sl)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		stat := validateMainStat(mana)

		person.mainStats[0] = byte(stat >> 2)
		rest := stat & uint16(0b0000000000000011)
		person.mainStats[1] = byte(rest << 6)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		stat := uint16(health)
		person.mainStats[1] = person.mainStats[1] | byte(stat>>8)
		rest := stat & uint16(0b0000000011111111)
		person.mainStats[2] = byte(rest)

		//val := uint32(person.mainStats[0])<<16 | uint32(person.mainStats[1])<<8 | uint32(person.mainStats[2])
		//fmt.Printf("%024b\n", val)
	}
}

func validateMainStat(stat int) uint16 {
	if stat > 1000 {
		stat = 1000
	}

	if stat < 0 {
		stat = 0
	}

	return uint16(stat)
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		stat := validateStat(respect)
		person.stats |= stat << 12
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		stat := validateStat(strength)
		person.stats |= stat << 8
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		stat := validateStat(experience)
		person.stats |= stat << 4
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		stat := validateStat(level)
		person.stats |= stat
	}
}

func validateStat(stat int) uint16 {
	if stat > 10 {
		stat = 10
	}

	if stat < 0 {
		stat = 0
	}

	return uint16(stat)
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.flags |= 1 << 4
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.flags |= 1 << 3
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.flags |= 1 << 2
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.flags |= byte(personType)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x, y, z int32
	gold    uint32

	//[number] -- bits count

	// [10]mana...gap...[10]health
	mainStats [3]byte

	//...[1]house[1]weapon[1]family[2]player-type
	// 000 0 0 0 00
	flags byte

	// [4]respect_[4]strength_[4]exp_[4]lvl
	stats uint16
	name  [42]byte
	/*
		name:
			42 байта
				Имя пользователя [0…42] символов латиницы
				символы латиницы это 65-122 ("A"-"z")

		x,y,z:
			по 4 байта
				оордината по оси X [-2_000_000_000…2_000_000_000] значений
				Координата по оси Y [-2_000_000_000…2_000_000_000] значений
				Координата по оси Z [-2_000_000_000…2_000_000_000] значений

		gold
			4 байта
				Золото [0…2_000_000_000] значений

		mainStats
			[0000000000_____0000000000____0000]
			3 байта:
				Магическая сила (мана) [0…1000] значений
				Здоровье [0…1000] значений

		stats:
			1 байт + побитовые маски (100%):
				Уважение [0…10] значений (0000????-1111????)
				Сила [0…10] значений     (????0000-????1111)

			1 байт + побитовые маски (100%):
				Опыт [0…10] значений      (0000????-1111????)
				Уровень [0…10] значений   (????0000-????1111)

		flags:
			1 байт + побитовые маски (100%)
				Есть ли у игрока дом [true/false] значения
				Есть ли у игрока оружие [true/false] значения
				Есть ли у игрока семья [true/false] значения
				Тип игрока [строитель/кузнец/воин] значения
	*/
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}

	for _, option := range options {
		option(&person)
	}

	return person
}

func (p *GamePerson) Name() string {
	return string(p.name[:])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	var res uint16
	res |= uint16(p.mainStats[0]) << 2
	rest := p.mainStats[1] & byte(0b11000000)
	res |= uint16(rest >> 6)

	// fmt.Printf("%016b\n", res)

	return int(res)
}

func (p *GamePerson) Health() int {
	//val := uint32(p.mainStats[0])<<16 | uint32(p.mainStats[1])<<8 | uint32(p.mainStats[2])
	//fmt.Printf("%024b\n", val)

	var res uint16
	leadingBits := p.mainStats[1] & byte(0b00000011)
	res |= uint16(leadingBits) << 8
	res |= uint16(p.mainStats[2])

	return int(res)
}

func (p *GamePerson) Respect() int {
	highBits := p.stats & uint16(0b1111000000000000)
	return int(highBits >> 12)
}

func (p *GamePerson) Strength() int {
	highBits := p.stats & uint16(0b0000111100000000)
	return int(highBits >> 8)
}

func (p *GamePerson) Experience() int {
	highBits := p.stats & uint16(0b0000000011110000)
	return int(highBits >> 4)
}

func (p *GamePerson) Level() int {
	highBits := p.stats & uint16(0b0000000000001111)
	return int(highBits)
}

func (p *GamePerson) HasHouse() bool {
	mask := uint8(0b00010000)
	return p.flags&mask > 0
}

func (p *GamePerson) HasGun() bool {
	mask := uint8(0b00001000)
	return p.flags&mask > 0
}

func (p *GamePerson) HasFamilty() bool {
	mask := uint8(0b00000100)
	return p.flags&mask > 0
}

func (p *GamePerson) Type() int {
	mask := uint8(0b00000011)
	return int(p.flags & mask)
}

func main() {
	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)

	fmt.Println(person.Respect())
}
