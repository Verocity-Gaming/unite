function score() {
    $.ajax({
        type: 'GET',
        dataType: 'json',
        url: "http://localhost:17069/http",
        timeout: 1000,
        success: function(data, status) {
            $('.purple').html(data.purple.value);
            $('.orange').html(data.orange.value);
            $('.self').html(data.self.value);
            $('.error').html('');
        },
        error: function(err) {
            $('.error').html("Failed to connect to Unite HUD server...");
            $('.purple').html("");
            $('.orange').html("");
            $('.self').html("");
        },
    });
};

$(document).ready(() => {
    setInterval(score, 1000);
});