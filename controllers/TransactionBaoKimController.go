package controllers
import (
	"net/http"
	"gopkg.in/unrolled/render.v1"
	"github.com/julienschmidt/httprouter"

)
type TransactionBaoKimController struct {
	AppController
	*render.Render
}
func (bk *TransactionBaoKimController) Index(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	
}

// URL checkout của baokim.vn
var  baokim_url = "https://www.baokim.vn/payment/customize_payment/order";

// Mã merchante site
var merchant_id = "647";	// Biến này được baokim.vn cung cấp khi bạn đăng ký merchant site

// Mật khẩu bảo mật
var secure_pass = "ae543c080ad91c23"; // Biến này được baokim.vn cung cấp khi bạn đăng ký merchant site

/**
 * Hàm xây dựng url chuyển đến BaoKim.vn thực hiện thanh toán, trong đó có tham số mã hóa (còn gọi là public key)
 * @param order_id				Mã đơn hàng
 * @param business 			Email tài khoản người bán
 * @param total_amount			Giá trị đơn hàng
 * @param shipping_fee			Phí vận chuyển
 * @param tax_fee				Thuế
 * @param order_description	Mô tả đơn hàng
 * @param url_success			Url trả về khi thanh toán thành công
 * @param url_cancel			Url trả về khi hủy thanh toán
 * @param url_detail			Url chi tiết đơn hàng
 * @return url cần tạo
 */
func (bk *TransactionBaoKimController) CreateRequestUrl(rw http.ResponseWriter, r *http.Request, p httprouter.Params){

}

/**
 * Hàm thực hiện xác minh tính chính xác thông tin trả về từ BaoKim.vn
 * @param _GET chứa tham số trả về trên url
 * @return true nếu thông tin là chính xác, false nếu thông tin không chính xác
 */
func (bk *TransactionBaoKimController)VerifyResponseUrl(rw http.ResponseWriter, r *http.Request, p httprouter.Params){

}
