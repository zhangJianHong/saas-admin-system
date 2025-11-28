#!/bin/bash

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查 Docker 是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    print_info "Docker 版本: $(docker --version)"
}

# 检查 Docker Compose 是否安装
check_docker_compose() {
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        print_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi
    if command -v docker-compose &> /dev/null; then
        print_info "Docker Compose 版本: $(docker-compose --version)"
        COMPOSE_CMD="docker-compose"
    else
        print_info "Docker Compose 版本: $(docker compose version)"
        COMPOSE_CMD="docker compose"
    fi
}

# 检查配置文件
check_config() {
    if [ ! -f "backend/configs/config.yaml" ]; then
        print_warn "配置文件不存在，是否从示例文件复制？(y/n)"
        read -r response
        if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
            if [ -f "backend/configs/config.yaml.example" ]; then
                cp backend/configs/config.yaml.example backend/configs/config.yaml
                print_info "已复制配置文件，请编辑 backend/configs/config.yaml 配置数据库连接信息"
                print_warn "请按任意键继续..."
                read -r
            else
                print_error "配置示例文件不存在"
                exit 1
            fi
        else
            print_error "请先创建配置文件 backend/configs/config.yaml"
            exit 1
        fi
    fi
}

# 检查环境变量文件
check_env() {
    if [ ! -f ".env" ]; then
        print_warn "环境变量文件 .env 不存在，是否从示例文件复制？(y/n)"
        read -r response
        if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
            if [ -f ".env.example" ]; then
                cp .env.example .env
                print_info "已复制环境变量文件，请编辑 .env 文件配置密码等敏感信息"
                print_warn "请按任意键继续..."
                read -r
            else
                print_error "环境变量示例文件不存在"
                exit 1
            fi
        else
            print_warn "将使用默认环境变量"
        fi
    fi
}

# 构建镜像
build_images() {
    print_info "开始构建 Docker 镜像..."
    $COMPOSE_CMD build --no-cache
    print_info "镜像构建完成"
}

# 启动服务
start_services() {
    print_info "启动服务..."
    $COMPOSE_CMD up -d
    print_info "服务启动完成"
}

# 查看服务状态
check_status() {
    print_info "服务状态:"
    $COMPOSE_CMD ps
}

# 查看日志
view_logs() {
    print_info "查看服务日志 (Ctrl+C 退出):"
    $COMPOSE_CMD logs -f
}

# 停止服务
stop_services() {
    print_info "停止服务..."
    $COMPOSE_CMD down
    print_info "服务已停止"
}

# 清理数据
clean_data() {
    print_warn "这将删除所有数据（包括数据库），是否继续？(yes/no)"
    read -r response
    if [[ "$response" == "yes" ]]; then
        print_info "停止并删除所有容器、网络和数据卷..."
        $COMPOSE_CMD down -v
        print_info "清理完成"
    else
        print_info "取消清理"
    fi
}

# 重启服务
restart_services() {
    print_info "重启服务..."
    $COMPOSE_CMD restart
    print_info "服务重启完成"
}

# 显示帮助信息
show_help() {
    cat << EOF
SaaS Monitor 部署脚本

用法: $0 [命令]

命令:
    deploy      - 完整部署（检查环境、构建镜像、启动服务）
    build       - 仅构建镜像
    start       - 启动服务
    stop        - 停止服务
    restart     - 重启服务
    status      - 查看服务状态
    logs        - 查看服务日志
    clean       - 停止服务并清理所有数据
    help        - 显示帮助信息

示例:
    $0 deploy   # 完整部署
    $0 logs     # 查看日志
    $0 status   # 查看状态

EOF
}

# 主函数
main() {
    case "${1:-deploy}" in
        deploy)
            print_info "开始部署 SaaS Monitor 系统..."
            check_docker
            check_docker_compose
            check_config
            check_env
            build_images
            start_services
            check_status
            print_info ""
            print_info "部署完成！"
            print_info "前端访问地址: http://localhost:3000"
            print_info "后端访问地址: http://localhost:8080"
            print_info "健康检查: http://localhost:8080/health"
            print_info ""
            print_info "默认登录账号: admin / admin123"
            print_info ""
            print_info "使用 '$0 logs' 查看日志"
            print_info "使用 '$0 status' 查看服务状态"
            ;;
        build)
            check_docker
            check_docker_compose
            build_images
            ;;
        start)
            check_docker
            check_docker_compose
            start_services
            check_status
            ;;
        stop)
            check_docker
            check_docker_compose
            stop_services
            ;;
        restart)
            check_docker
            check_docker_compose
            restart_services
            check_status
            ;;
        status)
            check_docker
            check_docker_compose
            check_status
            ;;
        logs)
            check_docker
            check_docker_compose
            view_logs
            ;;
        clean)
            check_docker
            check_docker_compose
            clean_data
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "未知命令: $1"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"
