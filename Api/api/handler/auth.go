package handler

import (
	"api/api/token"
	pb "api/genproto/auth"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user with email, password, and personal details
// @Security BearerAuth
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body pb.SignUpRequest true "User registration details"
// @Success 200 {object} pb.SignUpResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/signup [post]
func (h *Handler) SignUp(ctx *gin.Context) {
	var req pb.SignUpRequest
	var user pb.Users
	var reqq pb.EmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tkn := token.GenerateJWTToken(req.Email, "user")

	user.FirstName = req.LastName
	user.LastName = req.FirstName
	user.Email = req.Email
	reqq.Email = req.Email
	user.PhoneNumber = req.PhoneNumber
	user.Password = req.Password
	user.DateOfBirth = req.DateOfBirth
	user.AccessToken = tkn.AccessToken
	user.RefreshToken = tkn.RefreshToken

	reqq.Email = req.Email
	_, err := h.UserEmailSending(&reqq)
	if err != nil {
		panic(err)
	}

	response, err := h.Auth.SignUp(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)

}

// LogIn godoc
// @Summary Log in a user
// @Description Log in a user with email and password
// @Security BearerAuth
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body pb.LogInRequest true "User login details"
// @Success 200 {object} pb.LogInResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/login [post]
func (h *Handler) LogIn(ctx *gin.Context) {
	var req pb.LogInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.Auth.LogIn(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change password for a user
// @Security BearerAuth
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body pb.ChangePasswordRequest true "Change password details"
// @Success 200 {object} pb.InfoResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/change-password [post]
func (h *Handler) ChangePassword(ctx *gin.Context) {
	var req pb.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserId = "2bb835aa-5ad8-4d2e-b6db-08ca313c6ab9"

	response, err := h.Auth.ChangePassword(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ForgetPassword godoc
// @Summary Initiate password reset
// @Description Request a password reset email
// @Security BearerAuth
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body pb.ForgetPasswordRequest true "Request password reset"
// @Success 200 {object} pb.InfoResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/forget-password [post]
func (h *Handler) ForgetPassword(ctx *gin.Context) {
	var req pb.ForgetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.UserEmailSending((*pb.EmailRequest)(&req))
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, response) // Email sendend pleas confirm in Verify password
}

// ResetPassword godoc
// @Summary Reset user password
// @Description Reset a user's password using a temporary password
// @Security BearerAuth
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body pb.ResetPasswordRequest true "Reset password details"
// @Success 200 {object} pb.InfoResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/reset-password [post]
func (h *Handler) ResetPassword(ctx *gin.Context) {
	var req pb.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.CheckTheVerificationCode(ctx, &req)
	if err != nil {
		panic(err)
	}

	response, err := h.Auth.ResetPassword(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ChangeEmail godoc
// @Summary Change user email
// @Description Change the email address of a user
// @Security BearerAuth
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body pb.ChangeEmailRequest true "Change email details"
// @Success 200 {object} pb.InfoResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/change-email [post]
func (h *Handler) ChangeEmail(ctx *gin.Context) {
	var req pb.ChangeEmailRequest
	var email_req pb.EmailRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email_req.Email = req.NewEmail

	_, err := h.UserEmailSending(&email_req)
	if err != nil {
		panic(err)
	}

	response, err := h.Auth.ChangeEmail(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response) // You have to verify your new email
}

// VerifyEmail godoc
// @Summary Verify user email
// @Description Verify a user's email address with a verification code
// @Security BearerAuth
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body pb.VerifyEmailRequest true "Verify email details"
// @Success 200 {object} pb.InfoResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/verify-email [post]
func (h *Handler) VerifyEmail(ctx *gin.Context) {
	var req pb.VerifyEmailRequest
	var pass_req pb.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pass_req.Email = req.Email
	pass_req.VerificationCode = req.VerificationCode

	_, err := h.CheckTheVerificationCode(ctx, &pass_req)
	if err != nil {
		panic(err)
	}

	response, err := h.Auth.VerifyEmail(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// EnterEmail godoc
// @Summary Enter user email
// @Description Enter an email address for various purposes (e.g., account recovery)
// @Security BearerAuth
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body pb.EmailRequest true "Enter email details"
// @Success 200 {object} pb.InfoResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/enter-email [post]
func (h *Handler) EnterEmail(ctx *gin.Context) {
	var req pb.EmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.UserEmailSending(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ValidateToken godoc
// @Summary Validate an access token
// @Description Validate a user's access token
// @Security BearerAuth
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body pb.ValidateTokenRequest true "Token validation details"
// @Success 200 {object} pb.InfoResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/validateToken [post]
func (h *Handler) ValidateToken(c *gin.Context) {
	var req pb.ValidateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Auth.ValidateToken(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// RefreshToken godoc
// @Summary Refresh an access token
// @Description Refresh a user's access token using the refresh token
// @Security BearerAuth
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body pb.RefreshTokenRequest true "Token refresh details"
// @Success 200 {object} pb.InfoResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/refreshToken [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req pb.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Auth.RefreshToken(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
