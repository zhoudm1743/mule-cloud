package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/zhoudm1743/mule-cloud/internal/models"
	"github.com/zhoudm1743/mule-cloud/internal/repository"
	"github.com/zhoudm1743/mule-cloud/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 工序服务接口
type ProcessService interface {
	CreateProcess(ctx context.Context, req *models.CreateProcessRequest, userID primitive.ObjectID) (*models.Process, error)
	GetProcess(ctx context.Context, id primitive.ObjectID) (*models.Process, error)
	GetProcessByCode(ctx context.Context, code string) (*models.Process, error)
	UpdateProcess(ctx context.Context, id primitive.ObjectID, req *models.UpdateProcessRequest, userID primitive.ObjectID) (*models.Process, error)
	DeleteProcess(ctx context.Context, id primitive.ObjectID) error
	ListProcesses(ctx context.Context, req *models.ProcessListRequest) ([]*models.Process, int64, error)
	GetActiveProcesses(ctx context.Context) ([]*models.Process, error)
}

// 尺码服务接口
type SizeService interface {
	CreateSize(ctx context.Context, req *models.CreateSizeRequest, userID primitive.ObjectID) (*models.Size, error)
	GetSize(ctx context.Context, id primitive.ObjectID) (*models.Size, error)
	GetSizeByCode(ctx context.Context, code string) (*models.Size, error)
	UpdateSize(ctx context.Context, id primitive.ObjectID, req *models.UpdateSizeRequest, userID primitive.ObjectID) (*models.Size, error)
	DeleteSize(ctx context.Context, id primitive.ObjectID) error
	ListSizes(ctx context.Context, req *models.SizeListRequest) ([]*models.Size, int64, error)
	GetActiveSizes(ctx context.Context) ([]*models.Size, error)
}

// 颜色服务接口
type ColorService interface {
	CreateColor(ctx context.Context, req *models.CreateColorRequest, userID primitive.ObjectID) (*models.Color, error)
	GetColor(ctx context.Context, id primitive.ObjectID) (*models.Color, error)
	GetColorByCode(ctx context.Context, code string) (*models.Color, error)
	UpdateColor(ctx context.Context, id primitive.ObjectID, req *models.UpdateColorRequest, userID primitive.ObjectID) (*models.Color, error)
	DeleteColor(ctx context.Context, id primitive.ObjectID) error
	ListColors(ctx context.Context, req *models.ColorListRequest) ([]*models.Color, int64, error)
	GetActiveColors(ctx context.Context) ([]*models.Color, error)
}

// 客户服务接口
type CustomerService interface {
	CreateCustomer(ctx context.Context, req *models.CreateCustomerRequest, userID primitive.ObjectID) (*models.Customer, error)
	GetCustomer(ctx context.Context, id primitive.ObjectID) (*models.Customer, error)
	GetCustomerByCode(ctx context.Context, code string) (*models.Customer, error)
	UpdateCustomer(ctx context.Context, id primitive.ObjectID, req *models.UpdateCustomerRequest, userID primitive.ObjectID) (*models.Customer, error)
	DeleteCustomer(ctx context.Context, id primitive.ObjectID) error
	ListCustomers(ctx context.Context, req *models.CustomerListRequest) ([]*models.Customer, int64, error)
	GetActiveCustomers(ctx context.Context) ([]*models.Customer, error)
}

// 业务员服务接口
type SalespersonService interface {
	CreateSalesperson(ctx context.Context, req *models.CreateSalespersonRequest, userID primitive.ObjectID) (*models.Salesperson, error)
	GetSalesperson(ctx context.Context, id primitive.ObjectID) (*models.Salesperson, error)
	GetSalespersonByCode(ctx context.Context, code string) (*models.Salesperson, error)
	UpdateSalesperson(ctx context.Context, id primitive.ObjectID, req *models.UpdateSalespersonRequest, userID primitive.ObjectID) (*models.Salesperson, error)
	DeleteSalesperson(ctx context.Context, id primitive.ObjectID) error
	ListSalespersons(ctx context.Context, req *models.SalespersonListRequest) ([]*models.Salesperson, int64, error)
	GetActiveSalespersons(ctx context.Context) ([]*models.Salesperson, error)
}

// 工序服务实现
type processService struct {
	processRepo repository.ProcessRepository
	logger      logger.Logger
}

func NewProcessService(processRepo repository.ProcessRepository, logger logger.Logger) ProcessService {
	return &processService{
		processRepo: processRepo,
		logger:      logger,
	}
}

func (s *processService) CreateProcess(ctx context.Context, req *models.CreateProcessRequest, userID primitive.ObjectID) (*models.Process, error) {
	// 验证编码唯一性
	existing, err := s.processRepo.GetByCode(ctx, req.Code)
	if err != nil {
		s.logger.Error("Failed to check process code uniqueness", "code", req.Code, "error", err)
		return nil, fmt.Errorf("failed to check process code: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("process code '%s' already exists", req.Code)
	}

	// 创建工序
	process := &models.Process{
		Code:        strings.TrimSpace(req.Code),
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		UnitPrice:   req.UnitPrice,
		Category:    strings.TrimSpace(req.Category),
		SortOrder:   req.SortOrder,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	if err := s.processRepo.Create(ctx, process); err != nil {
		s.logger.Error("Failed to create process", "process", process, "error", err)
		return nil, fmt.Errorf("failed to create process: %w", err)
	}

	s.logger.Info("Process created successfully", "id", process.ID, "code", process.Code)
	return process, nil
}

func (s *processService) GetProcess(ctx context.Context, id primitive.ObjectID) (*models.Process, error) {
	process, err := s.processRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get process", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get process: %w", err)
	}
	if process == nil {
		return nil, fmt.Errorf("process not found")
	}

	return process, nil
}

func (s *processService) GetProcessByCode(ctx context.Context, code string) (*models.Process, error) {
	process, err := s.processRepo.GetByCode(ctx, code)
	if err != nil {
		s.logger.Error("Failed to get process by code", "code", code, "error", err)
		return nil, fmt.Errorf("failed to get process: %w", err)
	}
	if process == nil {
		return nil, fmt.Errorf("process not found")
	}

	return process, nil
}

func (s *processService) UpdateProcess(ctx context.Context, id primitive.ObjectID, req *models.UpdateProcessRequest, userID primitive.ObjectID) (*models.Process, error) {
	// 获取现有工序
	process, err := s.processRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get process: %w", err)
	}
	if process == nil {
		return nil, fmt.Errorf("process not found")
	}

	// 更新字段
	if req.Name != "" {
		process.Name = strings.TrimSpace(req.Name)
	}
	if req.Description != "" {
		process.Description = strings.TrimSpace(req.Description)
	}
	if req.UnitPrice > 0 {
		process.UnitPrice = req.UnitPrice
	}
	if req.Category != "" {
		process.Category = strings.TrimSpace(req.Category)
	}
	if req.IsActive != nil {
		process.IsActive = *req.IsActive
	}
	if req.SortOrder > 0 {
		process.SortOrder = req.SortOrder
	}
	process.UpdatedBy = userID

	if err := s.processRepo.Update(ctx, id, process); err != nil {
		s.logger.Error("Failed to update process", "id", id, "error", err)
		return nil, fmt.Errorf("failed to update process: %w", err)
	}

	s.logger.Info("Process updated successfully", "id", process.ID)
	return process, nil
}

func (s *processService) DeleteProcess(ctx context.Context, id primitive.ObjectID) error {
	// 检查工序是否存在
	process, err := s.processRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get process: %w", err)
	}
	if process == nil {
		return fmt.Errorf("process not found")
	}

	if err := s.processRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete process", "id", id, "error", err)
		return fmt.Errorf("failed to delete process: %w", err)
	}

	s.logger.Info("Process deleted successfully", "id", id)
	return nil
}

func (s *processService) ListProcesses(ctx context.Context, req *models.ProcessListRequest) ([]*models.Process, int64, error) {
	// 设置默认分页
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	processes, total, err := s.processRepo.List(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list processes", "error", err)
		return nil, 0, fmt.Errorf("failed to list processes: %w", err)
	}

	return processes, total, nil
}

func (s *processService) GetActiveProcesses(ctx context.Context) ([]*models.Process, error) {
	processes, err := s.processRepo.GetActive(ctx)
	if err != nil {
		s.logger.Error("Failed to get active processes", "error", err)
		return nil, fmt.Errorf("failed to get active processes: %w", err)
	}

	return processes, nil
}

// 尺码服务实现
type sizeService struct {
	sizeRepo repository.SizeRepository
	logger   logger.Logger
}

func NewSizeService(sizeRepo repository.SizeRepository, logger logger.Logger) SizeService {
	return &sizeService{
		sizeRepo: sizeRepo,
		logger:   logger,
	}
}

func (s *sizeService) CreateSize(ctx context.Context, req *models.CreateSizeRequest, userID primitive.ObjectID) (*models.Size, error) {
	// 验证编码唯一性
	existing, err := s.sizeRepo.GetByCode(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to check size code: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("size code '%s' already exists", req.Code)
	}

	size := &models.Size{
		Code:        strings.TrimSpace(req.Code),
		Name:        strings.TrimSpace(req.Name),
		Category:    strings.TrimSpace(req.Category),
		Description: strings.TrimSpace(req.Description),
		SortOrder:   req.SortOrder,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	if err := s.sizeRepo.Create(ctx, size); err != nil {
		return nil, fmt.Errorf("failed to create size: %w", err)
	}

	s.logger.Info("Size created successfully", "id", size.ID, "code", size.Code)
	return size, nil
}

func (s *sizeService) GetSize(ctx context.Context, id primitive.ObjectID) (*models.Size, error) {
	size, err := s.sizeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get size: %w", err)
	}
	if size == nil {
		return nil, fmt.Errorf("size not found")
	}
	return size, nil
}

func (s *sizeService) GetSizeByCode(ctx context.Context, code string) (*models.Size, error) {
	size, err := s.sizeRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get size: %w", err)
	}
	if size == nil {
		return nil, fmt.Errorf("size not found")
	}
	return size, nil
}

func (s *sizeService) UpdateSize(ctx context.Context, id primitive.ObjectID, req *models.UpdateSizeRequest, userID primitive.ObjectID) (*models.Size, error) {
	size, err := s.sizeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get size: %w", err)
	}
	if size == nil {
		return nil, fmt.Errorf("size not found")
	}

	// 更新字段
	if req.Name != "" {
		size.Name = strings.TrimSpace(req.Name)
	}
	if req.Category != "" {
		size.Category = strings.TrimSpace(req.Category)
	}
	if req.Description != "" {
		size.Description = strings.TrimSpace(req.Description)
	}
	if req.IsActive != nil {
		size.IsActive = *req.IsActive
	}
	if req.SortOrder > 0 {
		size.SortOrder = req.SortOrder
	}
	size.UpdatedBy = userID

	if err := s.sizeRepo.Update(ctx, id, size); err != nil {
		return nil, fmt.Errorf("failed to update size: %w", err)
	}

	s.logger.Info("Size updated successfully", "id", size.ID)
	return size, nil
}

func (s *sizeService) DeleteSize(ctx context.Context, id primitive.ObjectID) error {
	size, err := s.sizeRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get size: %w", err)
	}
	if size == nil {
		return fmt.Errorf("size not found")
	}

	if err := s.sizeRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete size: %w", err)
	}

	s.logger.Info("Size deleted successfully", "id", id)
	return nil
}

func (s *sizeService) ListSizes(ctx context.Context, req *models.SizeListRequest) ([]*models.Size, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	sizes, total, err := s.sizeRepo.List(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list sizes: %w", err)
	}

	return sizes, total, nil
}

func (s *sizeService) GetActiveSizes(ctx context.Context) ([]*models.Size, error) {
	sizes, err := s.sizeRepo.GetActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sizes: %w", err)
	}
	return sizes, nil
}

// 颜色服务实现
type colorService struct {
	colorRepo repository.ColorRepository
	logger    logger.Logger
}

func NewColorService(colorRepo repository.ColorRepository, logger logger.Logger) ColorService {
	return &colorService{
		colorRepo: colorRepo,
		logger:    logger,
	}
}

func (s *colorService) CreateColor(ctx context.Context, req *models.CreateColorRequest, userID primitive.ObjectID) (*models.Color, error) {
	existing, err := s.colorRepo.GetByCode(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to check color code: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("color code '%s' already exists", req.Code)
	}

	color := &models.Color{
		Code:        strings.TrimSpace(req.Code),
		Name:        strings.TrimSpace(req.Name),
		HexValue:    strings.TrimSpace(req.HexValue),
		RGBValue:    strings.TrimSpace(req.RGBValue),
		Category:    strings.TrimSpace(req.Category),
		Description: strings.TrimSpace(req.Description),
		SortOrder:   req.SortOrder,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	if err := s.colorRepo.Create(ctx, color); err != nil {
		return nil, fmt.Errorf("failed to create color: %w", err)
	}

	s.logger.Info("Color created successfully", "id", color.ID, "code", color.Code)
	return color, nil
}

func (s *colorService) GetColor(ctx context.Context, id primitive.ObjectID) (*models.Color, error) {
	color, err := s.colorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get color: %w", err)
	}
	if color == nil {
		return nil, fmt.Errorf("color not found")
	}
	return color, nil
}

func (s *colorService) GetColorByCode(ctx context.Context, code string) (*models.Color, error) {
	color, err := s.colorRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get color: %w", err)
	}
	if color == nil {
		return nil, fmt.Errorf("color not found")
	}
	return color, nil
}

func (s *colorService) UpdateColor(ctx context.Context, id primitive.ObjectID, req *models.UpdateColorRequest, userID primitive.ObjectID) (*models.Color, error) {
	color, err := s.colorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get color: %w", err)
	}
	if color == nil {
		return nil, fmt.Errorf("color not found")
	}

	if req.Name != "" {
		color.Name = strings.TrimSpace(req.Name)
	}
	if req.HexValue != "" {
		color.HexValue = strings.TrimSpace(req.HexValue)
	}
	if req.RGBValue != "" {
		color.RGBValue = strings.TrimSpace(req.RGBValue)
	}
	if req.Category != "" {
		color.Category = strings.TrimSpace(req.Category)
	}
	if req.Description != "" {
		color.Description = strings.TrimSpace(req.Description)
	}
	if req.IsActive != nil {
		color.IsActive = *req.IsActive
	}
	if req.SortOrder > 0 {
		color.SortOrder = req.SortOrder
	}
	color.UpdatedBy = userID

	if err := s.colorRepo.Update(ctx, id, color); err != nil {
		return nil, fmt.Errorf("failed to update color: %w", err)
	}

	s.logger.Info("Color updated successfully", "id", color.ID)
	return color, nil
}

func (s *colorService) DeleteColor(ctx context.Context, id primitive.ObjectID) error {
	color, err := s.colorRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get color: %w", err)
	}
	if color == nil {
		return fmt.Errorf("color not found")
	}

	if err := s.colorRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete color: %w", err)
	}

	s.logger.Info("Color deleted successfully", "id", id)
	return nil
}

func (s *colorService) ListColors(ctx context.Context, req *models.ColorListRequest) ([]*models.Color, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	colors, total, err := s.colorRepo.List(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list colors: %w", err)
	}
	return colors, total, nil
}

func (s *colorService) GetActiveColors(ctx context.Context) ([]*models.Color, error) {
	colors, err := s.colorRepo.GetActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active colors: %w", err)
	}
	return colors, nil
}

// 客户服务实现
type customerService struct {
	customerRepo repository.CustomerRepository
	logger       logger.Logger
}

func NewCustomerService(customerRepo repository.CustomerRepository, logger logger.Logger) CustomerService {
	return &customerService{
		customerRepo: customerRepo,
		logger:       logger,
	}
}

func (s *customerService) CreateCustomer(ctx context.Context, req *models.CreateCustomerRequest, userID primitive.ObjectID) (*models.Customer, error) {
	existing, err := s.customerRepo.GetByCode(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to check customer code: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("customer code '%s' already exists", req.Code)
	}

	customer := &models.Customer{
		Code:          strings.TrimSpace(req.Code),
		Name:          strings.TrimSpace(req.Name),
		ShortName:     strings.TrimSpace(req.ShortName),
		ContactPerson: strings.TrimSpace(req.ContactPerson),
		Phone:         strings.TrimSpace(req.Phone),
		Email:         strings.TrimSpace(req.Email),
		Address:       strings.TrimSpace(req.Address),
		TaxNumber:     strings.TrimSpace(req.TaxNumber),
		BankAccount:   strings.TrimSpace(req.BankAccount),
		PaymentTerms:  strings.TrimSpace(req.PaymentTerms),
		CreditLimit:   req.CreditLimit,
		CustomerType:  strings.TrimSpace(req.CustomerType),
		Region:        strings.TrimSpace(req.Region),
		Remarks:       strings.TrimSpace(req.Remarks),
		CreatedBy:     userID,
		UpdatedBy:     userID,
	}

	if err := s.customerRepo.Create(ctx, customer); err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	s.logger.Info("Customer created successfully", "id", customer.ID, "code", customer.Code)
	return customer, nil
}

func (s *customerService) GetCustomer(ctx context.Context, id primitive.ObjectID) (*models.Customer, error) {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	if customer == nil {
		return nil, fmt.Errorf("customer not found")
	}
	return customer, nil
}

func (s *customerService) GetCustomerByCode(ctx context.Context, code string) (*models.Customer, error) {
	customer, err := s.customerRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	if customer == nil {
		return nil, fmt.Errorf("customer not found")
	}
	return customer, nil
}

func (s *customerService) UpdateCustomer(ctx context.Context, id primitive.ObjectID, req *models.UpdateCustomerRequest, userID primitive.ObjectID) (*models.Customer, error) {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	if customer == nil {
		return nil, fmt.Errorf("customer not found")
	}

	// 更新字段
	if req.Name != "" {
		customer.Name = strings.TrimSpace(req.Name)
	}
	if req.ShortName != "" {
		customer.ShortName = strings.TrimSpace(req.ShortName)
	}
	if req.ContactPerson != "" {
		customer.ContactPerson = strings.TrimSpace(req.ContactPerson)
	}
	if req.Phone != "" {
		customer.Phone = strings.TrimSpace(req.Phone)
	}
	if req.Email != "" {
		customer.Email = strings.TrimSpace(req.Email)
	}
	if req.Address != "" {
		customer.Address = strings.TrimSpace(req.Address)
	}
	if req.TaxNumber != "" {
		customer.TaxNumber = strings.TrimSpace(req.TaxNumber)
	}
	if req.BankAccount != "" {
		customer.BankAccount = strings.TrimSpace(req.BankAccount)
	}
	if req.PaymentTerms != "" {
		customer.PaymentTerms = strings.TrimSpace(req.PaymentTerms)
	}
	if req.CreditLimit > 0 {
		customer.CreditLimit = req.CreditLimit
	}
	if req.Status != "" {
		customer.Status = strings.TrimSpace(req.Status)
	}
	if req.CustomerType != "" {
		customer.CustomerType = strings.TrimSpace(req.CustomerType)
	}
	if req.Region != "" {
		customer.Region = strings.TrimSpace(req.Region)
	}
	if req.Remarks != "" {
		customer.Remarks = strings.TrimSpace(req.Remarks)
	}
	if req.IsActive != nil {
		customer.IsActive = *req.IsActive
	}
	customer.UpdatedBy = userID

	if err := s.customerRepo.Update(ctx, id, customer); err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	s.logger.Info("Customer updated successfully", "id", customer.ID)
	return customer, nil
}

func (s *customerService) DeleteCustomer(ctx context.Context, id primitive.ObjectID) error {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}
	if customer == nil {
		return fmt.Errorf("customer not found")
	}

	if err := s.customerRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	s.logger.Info("Customer deleted successfully", "id", id)
	return nil
}

func (s *customerService) ListCustomers(ctx context.Context, req *models.CustomerListRequest) ([]*models.Customer, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	customers, total, err := s.customerRepo.List(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list customers: %w", err)
	}
	return customers, total, nil
}

func (s *customerService) GetActiveCustomers(ctx context.Context) ([]*models.Customer, error) {
	customers, err := s.customerRepo.GetActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active customers: %w", err)
	}
	return customers, nil
}

// 业务员服务实现
type salespersonService struct {
	salespersonRepo repository.SalespersonRepository
	logger          logger.Logger
}

func NewSalespersonService(salespersonRepo repository.SalespersonRepository, logger logger.Logger) SalespersonService {
	return &salespersonService{
		salespersonRepo: salespersonRepo,
		logger:          logger,
	}
}

func (s *salespersonService) CreateSalesperson(ctx context.Context, req *models.CreateSalespersonRequest, userID primitive.ObjectID) (*models.Salesperson, error) {
	existing, err := s.salespersonRepo.GetByCode(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to check salesperson code: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("salesperson code '%s' already exists", req.Code)
	}

	salesperson := &models.Salesperson{
		Code:       strings.TrimSpace(req.Code),
		Name:       strings.TrimSpace(req.Name),
		Phone:      strings.TrimSpace(req.Phone),
		Email:      strings.TrimSpace(req.Email),
		Department: strings.TrimSpace(req.Department),
		Position:   strings.TrimSpace(req.Position),
		HireDate:   req.HireDate,
		Region:     strings.TrimSpace(req.Region),
		Commission: req.Commission,
		Remarks:    strings.TrimSpace(req.Remarks),
		CreatedBy:  userID,
		UpdatedBy:  userID,
	}

	if err := s.salespersonRepo.Create(ctx, salesperson); err != nil {
		return nil, fmt.Errorf("failed to create salesperson: %w", err)
	}

	s.logger.Info("Salesperson created successfully", "id", salesperson.ID, "code", salesperson.Code)
	return salesperson, nil
}

func (s *salespersonService) GetSalesperson(ctx context.Context, id primitive.ObjectID) (*models.Salesperson, error) {
	salesperson, err := s.salespersonRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get salesperson: %w", err)
	}
	if salesperson == nil {
		return nil, fmt.Errorf("salesperson not found")
	}
	return salesperson, nil
}

func (s *salespersonService) GetSalespersonByCode(ctx context.Context, code string) (*models.Salesperson, error) {
	salesperson, err := s.salespersonRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get salesperson: %w", err)
	}
	if salesperson == nil {
		return nil, fmt.Errorf("salesperson not found")
	}
	return salesperson, nil
}

func (s *salespersonService) UpdateSalesperson(ctx context.Context, id primitive.ObjectID, req *models.UpdateSalespersonRequest, userID primitive.ObjectID) (*models.Salesperson, error) {
	salesperson, err := s.salespersonRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get salesperson: %w", err)
	}
	if salesperson == nil {
		return nil, fmt.Errorf("salesperson not found")
	}

	// 更新字段
	if req.Name != "" {
		salesperson.Name = strings.TrimSpace(req.Name)
	}
	if req.Phone != "" {
		salesperson.Phone = strings.TrimSpace(req.Phone)
	}
	if req.Email != "" {
		salesperson.Email = strings.TrimSpace(req.Email)
	}
	if req.Department != "" {
		salesperson.Department = strings.TrimSpace(req.Department)
	}
	if req.Position != "" {
		salesperson.Position = strings.TrimSpace(req.Position)
	}
	if !req.HireDate.IsZero() {
		salesperson.HireDate = req.HireDate
	}
	if req.Region != "" {
		salesperson.Region = strings.TrimSpace(req.Region)
	}
	if req.Commission > 0 {
		salesperson.Commission = req.Commission
	}
	if req.Status != "" {
		salesperson.Status = strings.TrimSpace(req.Status)
	}
	if req.Remarks != "" {
		salesperson.Remarks = strings.TrimSpace(req.Remarks)
	}
	if req.IsActive != nil {
		salesperson.IsActive = *req.IsActive
	}
	salesperson.UpdatedBy = userID

	if err := s.salespersonRepo.Update(ctx, id, salesperson); err != nil {
		return nil, fmt.Errorf("failed to update salesperson: %w", err)
	}

	s.logger.Info("Salesperson updated successfully", "id", salesperson.ID)
	return salesperson, nil
}

func (s *salespersonService) DeleteSalesperson(ctx context.Context, id primitive.ObjectID) error {
	salesperson, err := s.salespersonRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get salesperson: %w", err)
	}
	if salesperson == nil {
		return fmt.Errorf("salesperson not found")
	}

	if err := s.salespersonRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete salesperson: %w", err)
	}

	s.logger.Info("Salesperson deleted successfully", "id", id)
	return nil
}

func (s *salespersonService) ListSalespersons(ctx context.Context, req *models.SalespersonListRequest) ([]*models.Salesperson, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	salespersons, total, err := s.salespersonRepo.List(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list salespersons: %w", err)
	}
	return salespersons, total, nil
}

func (s *salespersonService) GetActiveSalespersons(ctx context.Context) ([]*models.Salesperson, error) {
	salespersons, err := s.salespersonRepo.GetActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active salespersons: %w", err)
	}
	return salespersons, nil
}
