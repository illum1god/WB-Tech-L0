document.addEventListener('DOMContentLoaded', function () {
    const getAllOrdersButton = document.getElementById('getAllOrders');
    const ordersListDiv = document.getElementById('ordersList');
    const searchOrderForm = document.getElementById('searchOrderForm');
    const orderDetailsDiv = document.getElementById('orderDetails');

    if (getAllOrdersButton) {
        getAllOrdersButton.addEventListener('click', function () {
            fetch('/api/order', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'X-CSRF-Token': getCSRFToken()
                }
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok ' + response.statusText);
                    }
                    return response.json();
                })
                .then(data => {
                    if (data.length === 0) {
                        ordersListDiv.innerHTML = '<p>No orders found.</p>';
                        return;
                    }

                    let html = '<ul>';
                    data.forEach(order => {
                        html += `<li class="order-item">
                                    <strong>Order UID:</strong> ${order.order_uid}<br>
                                    <strong>Track Number:</strong> ${order.track_number}<br>
                                    <a href="/order/${order.order_uid}">View Details</a>
                                 </li>`;
                    });
                    html += '</ul>';
                    ordersListDiv.innerHTML = html;
                })
                .catch(error => {
                    console.error('There has been a problem with your fetch operation:', error);
                    ordersListDiv.innerHTML = '<p>Error fetching orders.</p>';
                });
        });
    }

    if (searchOrderForm) {
        searchOrderForm.addEventListener('submit', function (e) {
            e.preventDefault();
            const uid = document.getElementById('uid').value.trim();
            if (!uid) {
                alert('Please enter Order UID.');
                return;
            }

            fetch(`/api/order/${uid}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'X-CSRF-Token': getCSRFToken()
                }
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok ' + response.statusText);
                    }
                    return response.json();
                })
                .then(data => {
                    displayOrderDetails(data);
                })
                .catch(error => {
                    console.error('There has been a problem with your fetch operation:', error);
                    orderDetailsDiv.innerHTML = '<p>Error fetching order details.</p>';
                });
        });
    }

    function displayOrderDetails(order) {
        let html = `<p><strong>Order UID:</strong> ${order.order_uid}</p>
                    <p><strong>Track Number:</strong> ${order.track_number}</p>
                    <p><strong>Entry:</strong> ${order.entry}</p>
                    <h3>Delivery</h3>
                    <p><strong>Name:</strong> ${order.delivery.name}</p>
                    <p><strong>Phone:</strong> ${order.delivery.phone}</p>
                    <p><strong>ZIP:</strong> ${order.delivery.zip}</p>
                    <p><strong>City:</strong> ${order.delivery.city}</p>
                    <p><strong>Address:</strong> ${order.delivery.address}</p>
                    <p><strong>Region:</strong> ${order.delivery.region}</p>
                    <p><strong>Email:</strong> ${order.delivery.email}</p>
                    <h3>Payment</h3>
                    <p><strong>Transaction:</strong> ${order.payment.transaction}</p>
                    <p><strong>Request ID:</strong> ${order.payment.request_id}</p>
                    <p><strong>Currency:</strong> ${order.payment.currency}</p>
                    <p><strong>Provider:</strong> ${order.payment.provider}</p>
                    <p><strong>Amount:</strong> ${order.payment.amount}</p>
                    <p><strong>Payment Date:</strong> ${new Date(order.payment.payment_dt * 1000).toLocaleString()}</p>
                    <p><strong>Bank:</strong> ${order.payment.bank}</p>
                    <p><strong>Delivery Cost:</strong> ${order.payment.delivery_cost}</p>
                    <p><strong>Goods Total:</strong> ${order.payment.goods_total}</p>
                    <p><strong>Custom Fee:</strong> ${order.payment.custom_fee}</p>
                    <h3>Items</h3>
                    <ul>`;
        order.items.forEach(item => {
            html += `<li>
                        <p><strong>Chrt ID:</strong> ${item.chrt_id}</p>
                        <p><strong>Track Number:</strong> ${item.track_number}</p>
                        <p><strong>Price:</strong> ${item.price}</p>
                        <p><strong>RID:</strong> ${item.rid}</p>
                        <p><strong>Name:</strong> ${item.name}</p>
                        <p><strong>Sale:</strong> ${item.sale}</p>
                        <p><strong>Size:</strong> ${item.size}</p>
                        <p><strong>Total Price:</strong> ${item.total_price}</p>
                        <p><strong>NM ID:</strong> ${item.nm_id}</p>
                        <p><strong>Brand:</strong> ${item.brand}</p>
                        <p><strong>Status:</strong> ${item.status}</p>
                    </li>`;
        });
        html += `</ul>`;
        orderDetailsDiv.innerHTML = html;
    }

    function getCSRFToken() {
        const csrfCookie = document.cookie.split('; ').find(row => row.startsWith('csrf_token='));
        if (csrfCookie) {
            return csrfCookie.split('=')[1];
        }
        return '';
    }
});