/***
 *
 * KLOUDS ORG
 *
 */

(function (module) {
    module.loader = function () {
        jQuery(".status").fadeOut();
        // will fade out the whole DIV that covers the website.
        jQuery(".preloader").delay(1000).fadeOut("slow");
    };
    module.initVideo = function () {
        $(".video-container").fitVids();
    };
    module.textEffect = function () {
        $(".texts > li").hide();
        setInterval(function () {
            $('.texts > li:first')
                .fadeOut(1000)
                .next()
                .fadeIn(1000)
                .end()
                .appendTo('.texts');
        }, 3000);

    };
    module.navPage = function () {
        $('.main-navigation').onePageNav({
            scrollThreshold: 0.2, // Adjust if Navigation highlights too early or too late
            filter: ':not(.external)',
            changeHash: true
        });
    };
    module.mainMediumNav = function () {
        if ($("#home").length > 0) {
            var top = (document.documentElement && document.documentElement.scrollTop) || document.body.scrollTop;
            if (top > 40) $('.sticky-navigation').stop().css({"top": '0',"position": 'fixed',"width":"100%"});
            else $('.sticky-navigation').stop().css({"position": 'relative'});
        }

    };
    module.mainMinNav = function () {
        var top = (document.documentElement && document.documentElement.scrollTop) || document.body.scrollTop;
        if (top > 40) $('.sticky-navigation').stop().animate({"top": '0',"position": 'fixed'});

        else $('.sticky-navigation').stop().css({"top": '-120',"position": 'fixed'});
    };
    module.mainSmallNav = function () {
        $('.main-navigation a').on('click', function () {
            $(".navbar-toggle").click();
        });
    };
    module.fullScreen = function () {
        $(".full-screen").css('min-height', $(window).height());
    };
    module.scrollTo = function () {
        var scrollAnimationTime = 1200,
            scrollAnimation = 'easeInOutExpo';
        $('a.scrollto').bind('click.smoothscroll', function (event) {
            event.preventDefault();
            var target = this.hash;
            $('html, body').stop().animate({
                'scrollTop': $(target).offset().top
            }, scrollAnimationTime, scrollAnimation, function () {
                window.location.hash = target;
            });
        });
    }
    module.init = function () {
        module.textEffect();
        module.navPage();
        module.initVideo();

        jQuery("#responsive_headline").fitText();

    };
})(Klouds = {});

/* =================================
 LOAD JS
 =================================== */

jQuery(window).load(function () {
    Klouds.loader();
})

jQuery(function () {
    //Klouds.mainMediumNav();
    Klouds.init();
    Klouds.fullScreen();
    Klouds.scrollTo();

    $(window).bind('resize', Klouds.fullScreen());
    wow = new WOW(
        {
            mobile: false
        });
    wow.init();
    $('#screenshots a').nivoLightbox({
        effect: 'fadeScale',
    });
    $("#feedbacks").owlCarousel({

        navigation: false, // Show next and prev buttons
        slideSpeed: 400,
        paginationSpeed: 800,
        autoPlay: 5000,
        singleItem: true
    });
    $('.expand-form').simpleexpand({
        'defaultTarget': '.expanded-contact-form'
    });


    $(window).stellar({
        horizontalScrolling: false
    });

});

$(window).scroll(function () {
    Klouds.mainMediumNav();
})


/* Bootstrap Internet Explorer 10 in Windows 8 and Windows Phone 8 FIX */

if (navigator.userAgent.match(/IEMobile\/10\.0/)) {
    var msViewportStyle = document.createElement('style')
    msViewportStyle.appendChild(
        document.createTextNode(
            '@-ms-viewport{width:auto!important}'
        )
    )
    document.querySelector('head').appendChild(msViewportStyle)
}
/* COLLAPSE NAVIGATION ON MOBILE AFTER CLICKING ON LINK*/

if (matchMedia('(max-width: 480px)').matches) {
    Klouds.mainSmallNav();
}

if (matchMedia('(min-width: 992px), (max-width: 767px)').matches) {
    Klouds.mainMediumNav();
}

if (matchMedia('(min-width: 768px) and (max-width: 991px)').matches) {
    Klouds.mainMinNav();
}

