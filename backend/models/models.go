package models

import (
	"time"

	"gorm.io/gorm"
)

type GormModel struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Occupant struct {
	gorm.Model
	Name                      string                  `json:"name" gorm:"not null"`
	Surname                   string                  `json:"surname" gorm:"not null"`
	Email                     string                  `json:"email"`
	Corporation               bool                    `json:"corporation" gorm:"default:false"`
	Phone                     string                  `json:"phone"`
	Iban                      string                  `json:"iban"`
	BirthDate                 time.Time               `json:"birth_date" gorm:"not null"`
	CivilityID                int                     `json:"civility_id" gorm:"not null"`
	Civility                  Civility                `gorm:"foreignKey:CivilityID"`
	DomicileAddressID         uint                    `json:"domicile_address_id" gorm:"not null"`
	DomicileAddress           Address                 `gorm:"foreignKey:DomicileAddressID"`
	DocumentReceivingMethodID int                     `json:"document_receiving_method_id" gorm:"default:0"`
	DocumentReceivingMethod   DocumentReceivingMethod `gorm:"foreignKey:DocumentReceivingMethodID"`
	ReminderDelay             int                     `json:"reminder_delay" gorm:"default:10"`
	ReminderReceivingMethodID int                     `json:"reminder_receiving_method_id" gorm:"default:0"`
	ReminderReceivingMethod   ReminderReceivingMethod `gorm:"foreignKey:ReminderReceivingMethodID"`
}

type Civility struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Type string `json:"type" gorm:"not null"`
}

type Address struct {
	gorm.Model
	Street     string `json:"street" gorm:"not null"`
	Complement string `json:"complement"`
	City       string `json:"city" gorm:"not null"`
	PostalCode string `json:"postal_code" gorm:"not null"`
	Country    string `json:"country" gorm:"not null"`
}

type DocumentReceivingMethod struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Type string `json:"type" gorm:"not null"`
}
type ReminderReceivingMethod struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Type string `json:"type" gorm:"not null"`
}

type OccupantPossessionOnProperty struct {
	OccupantID    uint         `json:"occupant_id" gorm:"not null"`
	Occupant      Occupant     `gorm:"foreignKey:OccupantID"`
	PropertyID    uint         `json:"property_id" gorm:"not null"`
	Property      Property     `gorm:"foreignKey:PropertyID"`
	Quota         float64      `json:"quota" gorm:"not null, default:0"`
	Administrator bool         `json:"administrator" gorm:"default:false"`
	OccupantType  uint         `json:"occupant_type" gorm:"not null"`
	OccupantTypes OccupantType `gorm:"foreignKey:OccupantType"`
}

type OccupantType struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Label string `json:"label" gorm:"not null"`
}

type Property struct {
	gorm.Model
	CondominiumID      uint             `json:"condominium_id" gorm:"not null"`
	Condominium        Condominium      `gorm:"foreignKey:CondominiumID"`
	AddressID          uint             `json:"address_id" gorm:"not null"`
	Address            Address          `gorm:"foreignKey:AddressID"`
	CadastralReference string           `json:"internal_reference" gorm:"not null"`
	PropertyType       uint             `json:"property_type" gorm:"not null"`
	PropertyTypes      PropertyType     `gorm:"foreignKey:PropertyType"`
	Floor              uint             `json:"floor" gorm:"not null"`
	Description        string           `json:"description"`
	Quota              float64          `json:"quota" gorm:"not null, default:0"`
	ElectricGazMeterID uint             `json:"electric_gaz_meter_id" gorm:"not null"`
	ElectricGazMeter   ElectricGazMeter `gorm:"foreignKey:ElectricGazMeterID"`
}

type ElectricGazMeter struct {
	gorm.Model
	Number       string `json:"number" gorm:"not null"`
	FtpImagePath string `json:"ftp_image_path" gorm:"not null"`
}

type PropertyType struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Label string `json:"label" gorm:"not null"`
}

type Condominium struct {
	gorm.Model
	Name               string   `json:"name" gorm:"not null"`
	AddressID          uint     `json:"address_id" gorm:"not null"`
	Address            Address  `gorm:"foreignKey:AddressID"`
	Description        string   `json:"description"`
	FtpBlueprintPath   string   `json:"ftp_blueprint_path" gorm:"not null"`
	LandRegistryNumber string   `json:"land_registry_number" gorm:"not null"`
	Prefix             string   `json:"prefix" gorm:"not null"`
	ConciergeID        uint     `json:"concierge_id" gorm:"not null"`
	Concierge          Occupant `gorm:"foreignKey:ConciergeID"`
}

type Exercice struct {
	gorm.Model
	CondominiumID uint        `json:"condominium_id" gorm:"not null"`
	Condominium   Condominium `gorm:"foreignKey:CondominiumID"`
	Date          time.Time   `json:"date" gorm:"not null"`
	Clotured      bool        `json:"clotured" gorm:"default:false"`
}

type Contract struct {
	gorm.Model
	SupplierID     uint        `json:"supplier_id" gorm:"not null"`
	Supplier       Supplier    `gorm:"foreignKey:SupplierID"`
	CondominiumID  uint        `json:"condominium_id" gorm:"not null"`
	Condominium    Condominium `gorm:"foreignKey:CondominiumID"`
	ContractNumber string      `json:"contract_number" gorm:"not null"`
	EndDate        time.Time   `json:"end_date" gorm:"not null"`
	StartDate      time.Time   `json:"start_date" gorm:"not null"`
	PriceInclVAT   float64     `json:"price_incl_vat" gorm:"not null"`
	PriceExclVAT   float64     `json:"price_excl_vat" gorm:"not null"`
	Description    string      `json:"description"`
}

type BankAccount struct {
	gorm.Model
	Iban          string          `json:"iban" gorm:"not null"`
	TypeID        uint            `json:"type_id" gorm:"not null"`
	Type          BankAccountType `gorm:"foreignKey:TypeID"`
	BankName      string          `json:"bank_name" gorm:"not null"`
	CondominiumID uint            `json:"condominium_id" gorm:"not null"`
	Condominium   Condominium     `gorm:"foreignKey:CondominiumID"`
	Balance       float64         `json:"balance" gorm:"not null, default:0"`
	Description   string          `json:"description"`
}

type BankAccountType struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Label string `json:"label" gorm:"not null"`
}

type AccountStatement struct {
	gorm.Model
	OperationDate      time.Time   `json:"operation_date" gorm:"not null"`
	ValueDate          time.Time   `json:"value_date" gorm:"not null"`
	Description        string      `json:"description" gorm:"not null"`
	Price              float64     `json:"price" gorm:"not null"`
	IbanID             uint        `json:"iban_id" gorm:"not null"`
	Iban               BankAccount `gorm:"foreignKey:IbanID"`
	InternalReference  string      `json:"internal_reference" gorm:"not null"`
	StatementReference string      `json:"statement_reference" gorm:"not null"`
	FtpFilePath        string      `json:"ftp_file_path" gorm:"not null"`
}

type AccountStatementOccupant struct {
	AccountStatementID uint             `json:"account_statement_id" gorm:"not null"`
	AccountStatement   AccountStatement `gorm:"foreignKey:AccountStatementID"`
	OccupantID         uint             `json:"occupant_id" gorm:"not null"`
	Occupant           Occupant         `gorm:"foreignKey:OccupantID"`
}

type GeneralAssembly struct {
	gorm.Model
	CondominiumID uint        `json:"condominium_id" gorm:"not null"`
	Condominium   Condominium `gorm:"foreignKey:CondominiumID"`
	FtpFilePath   string      `json:"ftp_file_path" gorm:"not null"`
	DateGA        time.Time   `json:"date_ga" gorm:"not null"`
	Clotured      bool        `json:"clotured" gorm:"default:false"`
}

type GAParticipation struct {
	OccupantID        uint            `json:"occupant_id" gorm:"not null"`
	Occupant          Occupant        `gorm:"foreignKey:OccupantID"`
	GeneralAssemblyID uint            `json:"general_assembly_id" gorm:"not null"`
	GeneralAssembly   GeneralAssembly `gorm:"foreignKey:GeneralAssemblyID"`
	Participation     bool            `json:"participation" gorm:"default:false"`
}

type CondoSupplier struct {
	SupplierID    uint        `json:"supplier_id" gorm:"not null"`
	Supplier      Supplier    `gorm:"foreignKey:SupplierID"`
	CondominiumID uint        `json:"condominium_id" gorm:"not null"`
	Condominium   Condominium `gorm:"foreignKey:CondominiumID"`
}

type Supplier struct {
	gorm.Model
	Name        string           `json:"name" gorm:"not null"`
	AddressID   uint             `json:"address_id" gorm:"not null"`
	Address     Address          `gorm:"foreignKey:AddressID"`
	Description string           `json:"description"`
	VATNumber   string           `json:"vat_number" gorm:"not null"`
	Phone       string           `json:"phone"`
	Email       string           `json:"email"`
	CategoryID  uint             `json:"category_id" gorm:"not null"`
	Category    SupplierCategory `gorm:"foreignKey:CategoryID"`
	EntryDate   time.Time        `json:"entry_date" gorm:"not null"`
	Iban        string           `json:"iban"`
}

type SupplierCategory struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Label string `json:"label" gorm:"not null"`
}

type Invoice struct {
	gorm.Model
	InvoiceType       bool        `json:"invoice_type" gorm:"not null default:true"`
	InvoiceNumber     string      `json:"invoice_number" gorm:"not null"`
	InternalReference string      `json:"internal_reference" gorm:"not null"`
	InvoiceLabel      string      `json:"invoice_label" gorm:"not null"`
	InvoiceDate       time.Time   `json:"invoice_date" gorm:"not null"`
	SupplierID        uint        `json:"supplier_id" gorm:"not null"`
	Supplier          Supplier    `gorm:"foreignKey:SupplierID"`
	CondominiumID     uint        `json:"condominium_id" gorm:"not null"`
	Condominium       Condominium `gorm:"foreignKey:CondominiumID"`
	PriceInclVAT      float64     `json:"price_incl_vat" gorm:"not null"`
	PriceExclVAT      float64     `json:"price_excl_vat" gorm:"not null"`
	InvoiceStatus     uint        `json:"invoice_status" gorm:"not null"`
	FtpFilePath       string      `json:"ftp_file_path" gorm:"not null"`
	ContractID        uint        `json:"contract_id" gorm:"not null"`
	Contract          Contract    `gorm:"foreignKey:ContractID"`
	ExerciceID        uint        `json:"exercice_id" gorm:"not null"`
	Exercice          Exercice    `gorm:"foreignKey:ExerciceID"`
}

type AccountStatementInvoice struct {
	AccountStatementID uint             `json:"account_statement_id" gorm:"not null"`
	AccountStatement   AccountStatement `gorm:"foreignKey:AccountStatementID"`
	InvoiceID          uint             `json:"invoice_id" gorm:"not null"`
	Invoice            Invoice          `gorm:"foreignKey:InvoiceID"`
}

type Remender struct {
	gorm.Model
	InvoicedID   uint      `json:"invoiced_id" gorm:"not null"`
	Date         time.Time `json:"date" gorm:"not null"`
	ReminderFees float64   `json:"reminder_fees" gorm:"not null"`
}

type AllocationKey struct {
	InvoiceID uint                  `json:"invoice_id" gorm:"not null"`
	Invoice   Invoice               `gorm:"foreignKey:InvoiceID"`
	KeyID     uint                  `json:"key_id" gorm:"not null"`
	Key       AllocationKeyTemplate `gorm:"foreignKey:KeyID"`
	PriceVAT  float64               `json:"price_vat" gorm:"not null"`
	Label     string                `json:"label" gorm:"not null"`
}

type AllocationKeyTemplate struct {
	ID    uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Code  int     `json:"code" gorm:"not null"`
	Label string  `json:"label" gorm:"not null"`
	VAT   float64 `json:"vat" gorm:"not null"`
}

type KeyProperty struct {
	KeyID      uint                  `json:"key_id" gorm:"not null"`
	Key        AllocationKeyTemplate `gorm:"foreignKey:KeyID"`
	PropertyID uint                  `json:"property_id" gorm:"not null"`
	Property   Property              `gorm:"foreignKey:PropertyID"`
}
