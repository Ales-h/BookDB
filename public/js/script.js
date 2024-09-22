document.querySelectorAll('.openFormButton').forEach(button => {
    button.addEventListener('click', function() {
        document.getElementById('modal').classList.remove('hidden');
    });
});


