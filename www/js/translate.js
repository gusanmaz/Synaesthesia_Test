function switchLang(obj) {
    var xhr = new XMLHttpRequest();
    var lang = obj.id.substring(0, 2)
    window.location.href = 'intro.html?lang=' + lang
}