package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		for i := 0; i < len(name); i++ {
			person.name[i] = name[i]
		}
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
		head := person.manaHouseGunFam & 0xFC00
		person.manaHouseGunFam = head | (uint16(mana) & 0x03ff)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		head := person.healthType & 0xFC00
		person.healthType = head | (uint16(health) & 0x03ff)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		head := person.respectPowerExpLvl & 0x0fff
		person.respectPowerExpLvl = head | (uint16(respect) << 12)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		head := person.respectPowerExpLvl & 0xf0ff
		person.respectPowerExpLvl = head | ((uint16(strength) << 8) & 0xf00)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		head := person.respectPowerExpLvl & 0xff0f
		person.respectPowerExpLvl = head | ((uint16(experience) << 4) & 0xf0)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		head := person.respectPowerExpLvl & 0xfff0
		person.respectPowerExpLvl = head | (uint16(level) & 0xf)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHouseGunFam |= 0x400
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHouseGunFam |= 0x800
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHouseGunFam |= 0x1000
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		head := person.healthType & 0xF3FF
		person.healthType = head | ((uint16(personType) & 0x3) << 10)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x, y, z            int32
	gold               uint32
	name               [42]byte
	manaHouseGunFam    uint16
	healthType         uint16
	respectPowerExpLvl uint16
}

func NewGamePerson(options ...Option) GamePerson {
	gamePerson := GamePerson{}
	for _, opt := range options {
		opt(&gamePerson)
	}
	return gamePerson
}

func (p *GamePerson) Name() string {
	sliceName := p.name[:]
	return unsafe.String(unsafe.SliceData(sliceName), len(sliceName))
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
	return int(p.manaHouseGunFam & 0x03ff)
}

func (p *GamePerson) Health() int {
	return int(p.healthType & 0x03ff)
}

func (p *GamePerson) Respect() int {
	return int((p.respectPowerExpLvl & 0xf000) >> 12)
}

func (p *GamePerson) Strength() int {
	return int((p.respectPowerExpLvl & 0x0f00) >> 8)
}

func (p *GamePerson) Experience() int {
	return int((p.respectPowerExpLvl & 0x00f0) >> 4)
}

func (p *GamePerson) Level() int {
	return int(p.respectPowerExpLvl & 0x000f)
}

func (p *GamePerson) HasHouse() bool {
	return (p.manaHouseGunFam & 0x400) == 0x400
}

func (p *GamePerson) HasGun() bool {
	return (p.manaHouseGunFam & 0x800) == 0x800
}

func (p *GamePerson) HasFamilty() bool {
	return (p.manaHouseGunFam & 0x1000) == 0x1000
}

func (p *GamePerson) Type() int {
	return int(p.healthType & 0xC00)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

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
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
