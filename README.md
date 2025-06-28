# Go Playground Plus

Una plataforma moderna para ejecutar código Go en un entorno seguro y controlado, con una interfaz web basada en Monaco Editor.

![Go Playground Plus Screenshot](https://images.unsplash.com/photo-1629654297299-c8506221ca97?auto=format&fit=crop&q=80&w=1000)

## Características Implementadas

### Arquitectura Modular

- **Diseño por Paquetes**: Código organizado en paquetes separados (`config`, `limiter`, `security`, `executor`, `handlers`, `logger`, `errors`)
- **Interfaces**: Uso de interfaces para mejorar modularidad y facilitar pruebas unitarias
- **Inyección de Dependencias**: Componentes desacoplados e inyectados donde se necesitan

### Seguridad

- **Validación de Código**: Análisis estático para detectar imports prohibidos usando el parser de Go
- **Sanitización de Entradas**: Validación estricta del código recibido
- **Límites de Ejecución**: Restricciones de tiempo y tamaño para el código ejecutado
- **Content Security Policy (CSP)**: Configuración robusta para prevenir XSS y otras vulnerabilidades
- **Headers de Seguridad**: X-Content-Type-Options, X-Frame-Options, etc.
- **CORS**: Configuración adecuada para solicitudes cross-origin

### Rendimiento

- **Rate Limiting**: Algoritmo de Token Bucket para control de tráfico eficiente
- **Pool de Buffers**: Uso de `sync.Pool` para reutilizar buffers y reducir la presión en el GC
- **Gestión de Recursos**: Cierre adecuado de recursos con `defer`
- **Timeout**: Control de tiempo máximo de ejecución para evitar bloqueos

### Logging y Manejo de Errores

- **Logging Estructurado**: Implementación con zap para logs eficientes y estructurados
- **Contexto en Errores**: Información adicional para facilitar debugging
- **Centralización**: Manejo centralizado de errores HTTP

### Despliegue

- **Containerización**: Configuración Docker optimizada
- **Docker Compose**: Orquestación de servicios
- **Volúmenes**: Montaje adecuado de archivos estáticos
- **Variables de Entorno**: Configuración externalizada

### Frontend

- **Monaco Editor**: Editor de código avanzado (mismo que usa VS Code)
- **Interfaz Moderna**: Diseño limpio y responsive con Tailwind CSS
- **Ejecución en Tiempo Real**: Visualización inmediata de resultados

## Getting Started

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/go-playground-plus.git
   ```

2. Navigate to the docker directory:

   ```bash
   cd go-playground-plus/docker
   ```

3. Start the application with Docker Compose:

   ```bash
   docker compose up -d
   ```

4. Open your browser and navigate to `http://localhost:8080`

## Uso

### Ejecución de Código Go

1. Escribe tu código Go en el editor
2. Haz clic en "Run" para ejecutar el código
3. Visualiza la salida en el panel de resultados
4. Los errores se mostrarán claramente en la consola

### Características del Editor

1. Resaltado de sintaxis
2. Autocompletado
3. Detección de errores
4. Soporte para temas claro/oscuro

## Stack Tecnológico

- **Backend**:

  - Go 1.21+
  - Arquitectura modular con interfaces
  - Logging estructurado (zap)
  - Manejo avanzado de errores

- **Frontend**:

  - React 18
  - TypeScript
  - Tailwind CSS
  - Monaco Editor
  - Lucide Icons

- **Infraestructura**:

  - Docker
  - Docker Compose
  - Volúmenes persistentes
  - Configuración por variables de entorno

## Contribuir

1. Haz un fork del repositorio

2. Crea tu rama de funcionalidad (`git checkout -b feature/NuevaFuncionalidad`)

3. Haz commit de tus cambios (`git commit -m 'Añadir nueva funcionalidad'`)

4. Sube los cambios a tu rama (`git push origin feature/NuevaFuncionalidad`)

5. Abre un Pull Request

## Licencia

Este proyecto está licenciado bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles.

## Agradecimientos

- [Go](https://go.dev/) - El lenguaje de programación Go
- [Monaco Editor](https://microsoft.github.io/monaco-editor/) - El editor de código que impulsa VS Code
- [Tailwind CSS](https://tailwindcss.com/) - Framework CSS basado en utilidades
- [Lucide](https://lucide.dev/) - Iconos consistentes y hermosos

## Contacto

Tu Nombre - [@tuusuario](https://twitter.com/tuusuario)

Project Link: [https://github.com/yourusername/go-playground-plus](https://github.com/yourusername/go-playground-plus)