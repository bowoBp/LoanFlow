package domians

import "time"

type (
	User struct {
		ID           uint      `gorm:"primaryKey;column:id" json:"id"`
		Email        string    `gorm:"unique;size:100;column:email" json:"email"`
		PasswordHash string    `gorm:"column:password_hash" json:"-"` // disembunyikan di JSON
		Role         string    `gorm:"size:50;column:role" json:"role"`
		Name         string    `gorm:"size:100;column:name" json:"name,omitempty"`
		Phone        string    `gorm:"size:20;column:phone" json:"phone,omitempty"`
		CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
		UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`

		// Relation (optional):
		// Satu user (role=BORROWER) bisa memiliki banyak loans
		BorrowerLoans []Loan `gorm:"foreignKey:BorrowerID;references:ID" json:"borrower_loans,omitempty"`

		// Jika user (role=STAFF) adalah "approver"
		LoanApprovalDetailsStaff []LoanApprovalDetail `gorm:"foreignKey:StaffID;references:ID" json:"approval_details_staff,omitempty"`

		// Jika user (role=STAFF) adalah "disburser"
		LoanDisbursementStaff []LoanDisbursementDetail `gorm:"foreignKey:StaffID;references:ID" json:"disbursement_staff,omitempty"`

		// Jika user (role=INVESTOR)
		LoanInvestors []LoanInvestor `gorm:"foreignKey:InvestorID;references:ID" json:"loan_investors,omitempty"`
	}

	RefreshToken struct {
		ID           uint      `gorm:"primaryKey;column:id" json:"id"`                            // Primary key
		UserID       uint      `gorm:"not null;column:user_id;index" json:"user_id"`              // FK to users.id
		RefreshToken string    `gorm:"not null;unique;column:refresh_token" json:"refresh_token"` // Unique refresh token
		ExpiresAt    time.Time `gorm:"not null;column:expires_at" json:"expires_at"`              // Expiry time of the token
		CreatedAt    time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`        // Time when the token was created
		UpdatedAt    time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`        // Time when the token was last updated

		// Relation with the User model
		User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	}
)
