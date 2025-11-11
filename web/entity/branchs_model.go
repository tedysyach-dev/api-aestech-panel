package entity

import "time"

type Branch struct {
	ID                   string     `gorm:"column:id;primaryKey;size:10"`
	Logo                 *string    `gorm:"column:logo;type:text"`
	NamaCabang           string     `gorm:"column:nama_cabang;size:150;not null"`
	Alamat               string     `gorm:"column:alamat;size:255;not null"`
	Kota                 string     `gorm:"column:kota;size:100;not null"`
	Kontak               string     `gorm:"column:kontak;size:25;not null"`
	Email                string     `gorm:"column:email;size:150;not null"`
	Koordinat            string     `gorm:"column:koordinat;size:255;not null"`
	Sipa                 string     `gorm:"column:sipa;size:255;not null"`
	IsPrivate            bool       `gorm:"column:is_private;not null;default:0"`
	Pettycash            int        `gorm:"column:pettycash;not null;default:0"`
	IP                   *string    `gorm:"column:ip;size:50"`
	KeyMachine           *string    `gorm:"column:key_machine;size:10"`
	IPMachine            *string    `gorm:"column:ip_machine;size:50"`
	BpomMode             bool       `gorm:"column:bpom_mode;default:0"`
	PPN                  int16      `gorm:"column:ppn;default:0"`
	Datetime             *string    `gorm:"column:datetime;size:200;default:'Asia/Jakarta'"`
	Upline               string     `gorm:"column:upline;size:10;not null"`
	IsManajemen          *bool      `gorm:"column:is_manajemen"`
	RoundPPN             *string    `gorm:"column:roundppn;size:10;default:'up'"`
	RegistDate           *time.Time `gorm:"column:regist_date;type:date"`
	ExpireDate           *time.Time `gorm:"column:expire_date;type:date"`
	IsPaid               *bool      `gorm:"column:is_paid;default:0"`
	Dev                  *bool      `gorm:"column:dev;default:0"`
	IsDelete             *bool      `gorm:"column:is_delete;default:0"`
	LastUpdateDashboard  *time.Time `gorm:"column:last_update_dashboard;type:date"`
	AvgGuest             *int       `gorm:"column:avg_guest;default:0"`
	AvgTransaction       *int       `gorm:"column:avg_transaction;default:0"`
	GuestCommentRate     *int       `gorm:"column:guest_comment_rate;default:0"`
	TopProduct           *string    `gorm:"column:top_product;type:text"`
	TopServices          *string    `gorm:"column:top_services;type:text"`
	TopProfAction        *string    `gorm:"column:top_prof_action;type:text"`
	GuestTotalByMonth    *int       `gorm:"column:guest_total_by_month;default:0"`
	TrxTotalByMonth      *int       `gorm:"column:trx_total_by_month;default:0"`
	ChartActivityByMonth *string    `gorm:"column:chart_activity_by_month;type:longtext"`
	ChartSalesByYear     *string    `gorm:"column:chart_sales_by_year;type:longtext"`
	RateReceptionist     *float64   `gorm:"column:rate_receptionist;default:0"`
	RateDoctor           *float64   `gorm:"column:rate_doctor;default:0"`
	RateBeautician       *float64   `gorm:"column:rate_beautician;default:0"`
	IDKlien              string     `gorm:"column:id_klien;size:255;not null"`
	NoWhatsapp           *string    `gorm:"column:no_whatsapp;size:50"`
	XenditID             *string    `gorm:"column:xendit_id;size:100"`
	WalletID             *string    `gorm:"column:wallet_id;size:100"`
	AccessID             *string    `gorm:"column:access_id;size:255"`
	AccessStatus         *bool      `gorm:"column:access_status;default:0"`
}

func (u *Branch) TableName() string {
	return "branchs"
}
