package constant

const (
	// Roles
	RoleAdmin    = "ADMIN"    // Untuk pengguna dengan akses penuh (web admin)
	RoleStaff    = "STAFF"    // Untuk staff operasional
	RoleBorrower = "BORROWER" // Untuk peminjam
	RoleInvestor = "INVESTOR" // Untuk investor

	// Role Descriptions (opsional untuk penjelasan lebih detail)
	RoleAdminDesc    = "Administrator with full access rights"
	RoleStaffDesc    = "Staff responsible for operational tasks"
	RoleBorrowerDesc = "User applying for loans"
	RoleInvestorDesc = "User investing in loans"
)
