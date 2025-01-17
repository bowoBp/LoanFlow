package domians

import "time"

type (
	Loan struct {
		ID                  uint      `gorm:"primaryKey;column:id" json:"id"`
		BorrowerID          uint      `gorm:"column:borrower_id" json:"borrower_id"`
		PrincipalAmount     float64   `gorm:"column:principal_amount" json:"principal_amount"`
		Rate                float64   `gorm:"column:rate" json:"rate"`
		ROI                 float64   `gorm:"column:roi" json:"roi"`
		State               string    `gorm:"size:20;column:state" json:"state"` // "proposed","approved","invested","disbursed"
		AgreementLetterLink string    `gorm:"column:agreement_letter_link" json:"agreement_letter_link,omitempty"`
		CreatedAt           time.Time `gorm:"column:created_at" json:"created_at"`
		UpdatedAt           time.Time `gorm:"column:updated_at" json:"updated_at"`

		// Relation back to user (borrower)
		Borrower User `gorm:"foreignKey:BorrowerID" json:"borrower,omitempty"`
		// One-to-Many to LoanApprovalDetail, LoanInvestor, LoanDisbursementDetail
		LoanApprovalDetails []LoanApprovalDetail     `gorm:"foreignKey:LoanID" json:"loan_approval_details,omitempty"`
		LoanInvestors       []LoanInvestor           `gorm:"foreignKey:LoanID" json:"loan_investors,omitempty"`
		LoanDisbursements   []LoanDisbursementDetail `gorm:"foreignKey:LoanID" json:"loan_disbursement_details,omitempty"`
	}

	LoanApprovalDetail struct {
		ID           uint      `gorm:"primaryKey;column:id" json:"id"`
		LoanID       uint      `gorm:"column:loan_id" json:"loan_id"`
		StaffID      uint      `gorm:"column:staff_id" json:"staff_id"`
		PhotoProof   string    `gorm:"column:photo_proof" json:"photo_proof,omitempty"`
		ApprovedDate time.Time `gorm:"column:approved_date" json:"approved_date,omitempty"`
		CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
		UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`

		// Relation back to Loan
		Loan Loan `gorm:"foreignKey:LoanID" json:"loan,omitempty"`
		// Relation back to Staff (User)
		Staff User `gorm:"foreignKey:StaffID" json:"staff,omitempty"`
	}

	LoanInvestor struct {
		ID             uint      `gorm:"primaryKey;column:id" json:"id"`
		LoanID         uint      `gorm:"column:loan_id" json:"loan_id"`
		InvestorID     uint      `gorm:"column:investor_id" json:"investor_id"`
		AmountInvested float64   `gorm:"column:amount_invested" json:"amount_invested"`
		CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
		UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`

		// Relation back to Loan
		Loan Loan `gorm:"foreignKey:LoanID" json:"loan,omitempty"`
		// Relation back to Investor (User)
		Investor User `gorm:"foreignKey:InvestorID" json:"investor,omitempty"`
	}

	LoanDisbursementDetail struct {
		ID                 uint      `gorm:"primaryKey;column:id" json:"id"`
		LoanID             uint      `gorm:"column:loan_id" json:"loan_id"`
		StaffID            uint      `gorm:"column:staff_id" json:"staff_id"`
		SignedAgreementDoc string    `gorm:"column:signed_agreement_doc" json:"signed_agreement_doc,omitempty"`
		DisbursedDate      time.Time `gorm:"column:disbursed_date" json:"disbursed_date,omitempty"`
		CreatedAt          time.Time `gorm:"column:created_at" json:"created_at"`
		UpdatedAt          time.Time `gorm:"column:updated_at" json:"updated_at"`

		// Relation back to Loan
		Loan Loan `gorm:"foreignKey:LoanID" json:"loan,omitempty"`
		// Relation back to Staff (User)
		Staff User `gorm:"foreignKey:StaffID" json:"staff,omitempty"`
	}

	LoanStateHistory struct {
		ID            uint      `gorm:"primaryKey;column:id" json:"id"`
		LoanID        uint      `gorm:"column:loan_id" json:"loan_id"`
		PreviousState string    `gorm:"column:previous_state" json:"previous_state"`
		NewState      string    `gorm:"column:new_state" json:"new_state"`
		ActionBy      uint      `gorm:"column:action_by" json:"action_by"`
		ActionAt      time.Time `gorm:"autoCreateTime;column:action_at" json:"action_at"`
		Remarks       string    `gorm:"column:remarks" json:"remarks"`
	}
)
