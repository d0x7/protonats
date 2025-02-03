plugins {
    id("java")
    id("maven-publish")
}

group = "com.xiamli"
version = "0.0.1-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
    testImplementation(platform("org.junit:junit-bom:5.10.0"))
    testImplementation("org.junit.jupiter:junit-jupiter")
    implementation("com.google.protobuf:protobuf-java:4.29.3")
}

tasks.test {
    useJUnitPlatform()
}

sourceSets {
    main {
        java {
            srcDir("src/main/java")
            srcDir("src/main/protobuf")
        }
    }
}

publishing {
    repositories {
        maven {
            name = "GitHubPackages"
            url = uri("https://maven.pkg.github.com/d0x7/protonats")
            credentials {
                username = project.findProperty("gpr.user") as String? ?: System.getenv("GH_USERNAME")
                password = project.findProperty("gpr.key") as String? ?: System.getenv("GH_TOKEN")
            }
        }
    }
    publications {
        register<MavenPublication>("gpr") {
            from(components["java"])
            artifactId = "protonats"
        }
    }
}